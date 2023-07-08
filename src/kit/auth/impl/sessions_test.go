package impl

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type sessionsTestSuite struct {
	kitTestSuite.Suite
	userRepository *mocks.UserRepository
	sessionStorage *mocks.SessionStorage
	authCfg        *auth.Config
	svc            auth.SessionsService
	logger         func() log.CLogger
	authorize      *mocks.AuthorizeSession
}

func (s *sessionsTestSuite) SetupSuite() {
	s.logger = func() log.CLogger {
		return log.L(log.Init(&log.Config{Level: log.InfoLevel}))
	}
	s.Suite.Init(s.logger)
	s.authCfg = &auth.Config{
		AccessToken: &auth.TokenConfig{
			Secret:              "123",
			ExpirationPeriodSec: 60,
		},
		RefreshToken: &auth.TokenConfig{
			Secret:              "321",
			ExpirationPeriodSec: 60,
		},
	}
}

func (s *sessionsTestSuite) SetupTest() {
	s.userRepository = &mocks.UserRepository{}
	s.sessionStorage = &mocks.SessionStorage{}
	s.authorize = &mocks.AuthorizeSession{}
	s.svc = NewSessionsService(s.logger, s.userRepository, s.sessionStorage, s.authorize)
	s.svc.Init(s.authCfg)
}

func TestSessionsSuite(t *testing.T) {
	suite.Run(t, new(sessionsTestSuite))
}

func (s *sessionsTestSuite) getUser() *auth.User {
	now := kit.Now()
	return &auth.User{
		Id:          "1",
		Username:    "user@test.com",
		Password:    "111",
		FirstName:   "user",
		LastName:    "user",
		ActivatedAt: &now,
		Groups:      []string{"group"},
	}
}

func (s *sessionsTestSuite) Test_LoginPassword_WhenUserNotFound_Fail() {
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(nil, nil)
	rq := &auth.LoginRequest{}
	_, _, err := s.svc.LoginPassword(s.Ctx, rq)
	s.AssertAppErr(err, auth.ErrCodeUserNotFound)
}

func (s *sessionsTestSuite) Test_LoginPassword_WhenUserNotActive_Fail() {
	usr := s.getUser()
	usr.ActivatedAt = nil
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(usr, nil)
	rq := &auth.LoginRequest{
		Username: usr.Username,
	}
	_, _, err := s.svc.LoginPassword(s.Ctx, rq)
	s.AssertAppErr(err, auth.ErrCodeUserNotActive)
}

func (s *sessionsTestSuite) Test_LoginPassword_WhenUserLocked_Fail() {
	usr := s.getUser()
	now := kit.Now()
	usr.LockedAt = &now
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(usr, nil)
	rq := &auth.LoginRequest{
		Username: usr.Username,
	}
	_, _, err := s.svc.LoginPassword(s.Ctx, rq)
	s.AssertAppErr(err, auth.ErrCodeUserLocked)
}

func (s *sessionsTestSuite) assertToken(token *auth.SessionToken) {
	svcImpl := s.svc.(*serviceImpl)
	_, _, err := svcImpl.verifyJwtToken(s.Ctx, token.AccessToken, s.authCfg.AccessToken.Secret)
	s.Nil(err)
	_, _, err = svcImpl.verifyJwtToken(s.Ctx, token.RefreshToken, s.authCfg.RefreshToken.Secret)
	s.Nil(err)
}

func (s *sessionsTestSuite) Test_LoginPassword_CreateSession_Ok() {
	password := "123456"
	usr := s.getUser()
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usr.Password = string(bytes)
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(usr, nil)
	roles := []string{"role"}
	s.authorize.On("GetRolesForGroups", s.Ctx, usr.Groups).Return(roles, nil)
	var actualSession *auth.Session
	s.sessionStorage.On("CreateSession", s.Ctx, mock.AnythingOfType("*auth.Session")).
		Run(func(args mock.Arguments) {
			actualSession = args.Get(1).(*auth.Session)
		}).Return(nil)
	rq := &auth.LoginRequest{
		Username: usr.Username,
		Password: password,
	}
	_, token, err := s.svc.LoginPassword(s.Ctx, rq)
	s.Nil(err)
	s.NotEmpty(actualSession)
	s.Equal(rq.Username, actualSession.Username)
	s.Equal(roles, actualSession.Roles)
	s.Equal(usr.Id, actualSession.UserId)
	s.NotEmpty(actualSession.LoginAt)
	s.NotEmpty(actualSession.LastActivityAt)
	s.NotEmpty(token)
	s.assertToken(token)
}

func (s *sessionsTestSuite) Test_Logout_NoActiveSessions_Ok() {
	session := &auth.Session{
		Id:     kit.NewId(),
		UserId: kit.NewId(),
	}
	user := s.getUser()
	s.sessionStorage.On("Get", s.Ctx, session.Id).Return(session, nil)
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(user, nil)
	s.sessionStorage.On("GetByUser", s.Ctx, session.UserId).Return(nil, nil)
	s.sessionStorage.On("Logout", s.Ctx, session.Id, mock.AnythingOfType("time.Time")).Return(nil)
	err := s.svc.Logout(s.Ctx, session.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.AssertNumberOfCalls(&s.sessionStorage.Mock, "Logout", 1)
}

func (s *sessionsTestSuite) Test_Logout_OtherActiveSessions_Ok() {
	session := &auth.Session{
		Id:     kit.NewId(),
		UserId: kit.NewId(),
	}
	otherSession := &auth.Session{
		Id:     kit.NewId(),
		UserId: session.UserId,
	}
	s.sessionStorage.On("Get", s.Ctx, session.Id).Return(session, nil)
	s.sessionStorage.On("Get", s.Ctx, otherSession.Id).Return(session, nil)
	onceCalled := false
	onceCalledFn := func() []*auth.Session {
		if !onceCalled {
			onceCalled = true
			return []*auth.Session{otherSession}
		}
		return nil
	}
	sessStorageMock := s.sessionStorage.On("GetByUser", s.Ctx, session.UserId)
	sessStorageMock.Run(func(args mock.Arguments) {
		sessStorageMock.ReturnArguments = mock.Arguments{onceCalledFn(), nil}
	})
	s.sessionStorage.On("Logout", s.Ctx, session.Id, mock.AnythingOfType("time.Time")).Return(nil)
	s.sessionStorage.On("Logout", s.Ctx, otherSession.Id, mock.AnythingOfType("time.Time")).Return(nil)
	err := s.svc.Logout(s.Ctx, session.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.AssertNumberOfCalls(&s.sessionStorage.Mock, "Logout", 2)
}

func (s *sessionsTestSuite) Test_RefreshToken_WhenValid_Success() {
	password := "123456"
	usr := s.getUser()
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usr.Password = string(bytes)
	rq := &auth.LoginRequest{
		Username: usr.Username,
		Password: password,
	}
	s.userRepository.On("GetByUsername", s.Ctx, mock.AnythingOfType("string")).Return(usr, nil)
	s.sessionStorage.On("CheckRefreshToken", s.Ctx, mock.AnythingOfType("string")).Return(true, nil)
	roles := []string{"role"}
	s.authorize.On("GetRolesForGroups", s.Ctx, usr.Groups).Return(roles, nil)
	s.sessionStorage.On("CreateSession", s.Ctx, mock.AnythingOfType("*auth.Session")).Return(nil)
	session, token, err := s.svc.LoginPassword(s.Ctx, rq)
	if err != nil {
		s.Fatal(err)
	}
	s.sessionStorage.On("Get", s.Ctx, session.Id).Return(session, nil)

	refreshedToken, err := s.svc.RefreshToken(s.Ctx, token.RefreshToken)
	if err != nil {
		s.Fatal(err)
	}
	s.assertToken(refreshedToken)
}

func (s *sessionsTestSuite) Test_WhenCreateTokensTwoTimesInARow_GetDifferentTokens_Ok() {
	session := &auth.Session{
		Id:       "1",
		UserId:   "1",
		Username: "1",
	}
	svc := s.svc.(*serviceImpl)
	tkn1, err := svc.createJwtToken(s.Ctx, session)
	if err != nil {
		s.Fatal(err)
	}
	tkn2, err := svc.createJwtToken(s.Ctx, session)
	if err != nil {
		s.Fatal(err)
	}
	s.NotEqual(tkn1.AccessToken, tkn2.AccessToken)
	s.NotEqual(tkn1.RefreshToken, tkn2.RefreshToken)
}
