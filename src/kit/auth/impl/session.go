package impl

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type serviceImpl struct {
	userRepository auth.UserRepository
	sessionStorage auth.SessionStorage
	authorize      auth.AuthorizeSession
	logger         log.CLoggerFunc
	cfg            *auth.Config
}

func NewSessionsService(logger log.CLoggerFunc, userRepository auth.UserRepository, sessionStorage auth.SessionStorage, authorize auth.AuthorizeSession) auth.SessionsService {
	return &serviceImpl{
		userRepository: userRepository,
		sessionStorage: sessionStorage,
		authorize:      authorize,
		logger:         logger,
	}
}

func (s *serviceImpl) l() log.CLogger {
	return s.logger().Cmp("sessions-svc")
}

func (s *serviceImpl) Init(cfg *auth.Config) {
	s.cfg = cfg
}

func (s *serviceImpl) createJwtToken(ctx context.Context, sess *auth.Session) (*auth.SessionToken, error) {
	s.l().C(ctx).Mth("create-jwt").F(log.FF{"uid": sess.UserId}).Dbg()

	st := &auth.SessionToken{
		SessionId: sess.Id,
	}

	now := kit.Now()

	// access token
	atExpireAt := now.Add(time.Second * time.Duration(s.cfg.AccessToken.ExpirationPeriodSec))
	atClaims := jwt.MapClaims{}
	atClaims["tid"] = kit.NewId()
	atClaims["exp"] = atExpireAt.Unix()
	atClaims["sid"] = sess.Id
	atClaims["uid"] = sess.UserId
	atClaims["un"] = sess.Username
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atStr, err := at.SignedString([]byte(s.cfg.AccessToken.Secret))
	if err != nil {
		return nil, auth.ErrAccessTokenCreation(err, ctx)
	}

	st.AccessToken = atStr
	st.AccessTokenExpiresAt = atExpireAt

	// refresh token
	rtExpireAt := now.Add(time.Second * time.Duration(s.cfg.RefreshToken.ExpirationPeriodSec))
	rtClaims := jwt.MapClaims{}
	rtClaims["tid"] = kit.NewId()
	rtClaims["exp"] = rtExpireAt.Unix()
	rtClaims["sid"] = sess.Id
	rtClaims["uid"] = sess.UserId
	rtClaims["un"] = sess.Username
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtStr, err := rt.SignedString([]byte(s.cfg.RefreshToken.Secret))
	if err != nil {
		return nil, auth.ErrAccessTokenCreation(err, ctx)
	}
	st.RefreshToken = rtStr
	st.RefreshTokenExpiresAt = rtExpireAt

	return st, nil
}

func (s *serviceImpl) createSession(ctx context.Context, usr *auth.User) (*auth.Session, *auth.SessionToken, error) {
	l := s.l().C(ctx).Mth("create-session").F(log.FF{"user": usr.Username}).Trc()

	// get roles by usr groups
	roles, err := s.authorize.GetRolesForGroups(ctx, usr.Groups)
	if err != nil {
		return nil, nil, err
	}

	// create a new session
	now := kit.Now()
	session := &auth.Session{
		Id:             kit.NewId(),
		UserId:         usr.Id,
		Username:       usr.Username,
		LoginAt:        now,
		LastActivityAt: now,
		Roles:          roles,
	}

	// create JWT
	token, err := s.createJwtToken(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	// save session to store
	if err := s.sessionStorage.CreateSession(ctx, session); err != nil {
		return nil, nil, err
	}

	l.F(log.FF{"sid": session.Id}).Dbg("ok")

	return session, token, nil
}

func (s *serviceImpl) checkUserPassword(ctx context.Context, rq *auth.LoginRequest, storedHash string) error {
	s.l().C(ctx).Mth("check-user-password").F(log.FF{"user": rq.Username}).Trc()
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(rq.Password))
	if err != nil {
		return auth.ErrSessionPasswordValidation(ctx)
	}
	return nil
}

func (s *serviceImpl) LoginPassword(ctx context.Context, rq *auth.LoginRequest) (*auth.Session, *auth.SessionToken, error) {
	s.l().C(ctx).Mth("login-password").F(log.FF{"user": rq.Username}).Trc()

	// get user by username
	usr, err := s.userRepository.GetByUsername(ctx, rq.Username)
	if err != nil {
		return nil, nil, err
	}
	if usr == nil {
		return nil, nil, auth.ErrUserNotFound(ctx, rq.Username)
	}
	// check user status
	if usr.ActivatedAt == nil {
		return nil, nil, auth.ErrUserNotActive(ctx, rq.Username)
	}
	if usr.LockedAt != nil {
		return nil, nil, auth.ErrUserLocked(ctx, rq.Username)
	}

	// check password
	err = s.checkUserPassword(ctx, rq, usr.Password)
	if err != nil {
		return nil, nil, err
	}
	return s.createSession(ctx, usr)
}

func (s *serviceImpl) Logout(ctx context.Context, sid string) error {
	l := s.l().C(ctx).Mth("logout").F(log.FF{"sid": sid}).Trc()

	// find user sessions
	ss, err := s.sessionStorage.Get(ctx, sid)
	if err != nil {
		return err
	}
	if ss == nil {
		l.Warn("no sessions found")
		return nil
	}
	l.F(log.FF{"uid": ss.UserId})

	// check session is already logged out
	if ss.LogoutAt != nil {
		return auth.ErrSessionLoggedOut(ctx)
	}

	// mark session as logged out
	if err := s.sessionStorage.Logout(ctx, sid, kit.Now()); err != nil {
		return err
	}
	l.Dbg("logged out")

	// logout all user's sessions
	err = s.logoutUser(ctx, ss.UserId)
	if err != nil {
		return err
	}

	l.Trc("ok")
	return nil
}

func (s *serviceImpl) verifyJwtToken(ctx context.Context, tokenStr string, secret string) (*jwt.Token, jwt.MapClaims, error) {

	// parse JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrSessionAuthWrongSigningMethod(ctx)
		}
		return []byte(secret), nil
	})
	if err != nil {
		if jwtErr, ok := err.(*jwt.ValidationError); ok {
			if jwtErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, nil, auth.ErrSessionAuthTokenExpired(ctx)
			}
		}
		return nil, nil, auth.ErrSessionAuthTokenInvalid(ctx)
	}
	if !token.Valid {
		return nil, nil, auth.ErrSessionAuthTokenInvalid(ctx)
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, nil, auth.ErrSessionAuthTokenClaimsInvalid(ctx)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, auth.ErrSessionAuthTokenClaimsInvalid(ctx)
	}

	return token, claims, nil
}

func (s *serviceImpl) AuthSession(ctx context.Context, accessToken string) (*auth.Session, error) {
	l := s.l().C(ctx).Mth("auth").Trc()

	// verify JWT token
	_, claims, err := s.verifyJwtToken(ctx, accessToken, s.cfg.AccessToken.Secret)
	if err != nil {
		return nil, err
	}

	// extract SID from claims
	sid := claims["sid"].(string)
	l.F(log.FF{"sid": sid})

	// get token by sid
	session, err := s.sessionStorage.Get(ctx, sid)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, auth.ErrSessionTokenInvalid(ctx)
	}

	l.F(log.FF{"uid": session.UserId})

	// check session is logged out
	if session.LogoutAt != nil {
		return nil, auth.ErrSessionLoggedOut(ctx)
	}

	l.Trc("ok")

	// update session's last activity asynchronously
	// we have to invent another way of updating to avoid too many DB hits (maybe periodic async cron)
	//go func() {
	//	if err := s.sessionStorage.UpdateLastActivity(ctx, session.Id, kit.Now()); err != nil {
	//		s.l().Mth("session-lastactivity").E(err).Err()
	//	}
	//}()

	return session, nil
}

func (s *serviceImpl) RefreshToken(ctx context.Context, refreshToken string) (*auth.SessionToken, error) {

	l := s.l().C(ctx).Mth("refresh-token").Dbg()

	// verify JWT token
	_, claims, err := s.verifyJwtToken(ctx, refreshToken, s.cfg.RefreshToken.Secret)
	if err != nil {
		return nil, err
	}

	// extract SID from claims
	sid := claims["sid"].(string)

	l.F(log.FF{"sid": sid})

	// get session
	session, err := s.sessionStorage.Get(ctx, sid)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, auth.ErrSessionNoSessionFound(ctx)
	}

	// check session is logged out
	if session.LogoutAt != nil {
		return nil, auth.ErrSessionLoggedOut(ctx)
	}

	// issue a new access token
	token, err := s.createJwtToken(ctx, session)
	if err != nil {
		return nil, err
	}

	l.Dbg("ok")

	return token, nil
}

func (s *serviceImpl) Get(ctx context.Context, sid string) (*auth.Session, error) {
	s.l().C(ctx).Mth("get").F(log.FF{"sid": sid}).Dbg()
	return s.sessionStorage.Get(ctx, sid)
}

func (s *serviceImpl) GetByUser(ctx context.Context, userId string) ([]*auth.Session, error) {
	s.l().C(ctx).Mth("get-by-user").F(log.FF{"userId": userId}).Trc()
	return s.sessionStorage.GetByUser(ctx, userId)
}

func (s *serviceImpl) logoutUser(ctx context.Context, userId string) error {
	s.l().C(ctx).Mth("logout-user").F(log.FF{"uid": userId}).Trc()
	// get all user's sessions
	sessions, err := s.sessionStorage.GetByUser(ctx, userId)
	if err != nil {
		return err
	}
	// logout all sessions
	for _, ss := range sessions {
		err = s.Logout(ctx, ss.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
