package auth

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"golang.org/x/crypto/bcrypt"
	"unicode/utf8"
)

const (
	PasswordMinLen = 6
)

type userSvcImpl struct {
	storage domain.UserStorage
}

func NewUserService(storage domain.UserStorage) domain.UserService {
	return &userSvcImpl{
		storage: storage,
	}
}

func (u *userSvcImpl) l() log.CLogger {
	return service.L().Cmp("users-domain-svc")
}

func (u *userSvcImpl) getPwdHash(ctx context.Context, password string) (string, string, error) {

	if password == "" {
		return "", "", errors.ErrAuthPwdEmpty(ctx)
	}
	// check password policy
	if utf8.RuneCountInString(password) < PasswordMinLen {
		return "", "", errors.ErrAuthPwdPolicy(ctx)
	}

	// hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", errors.ErrUserPasswordHashGenerate(err, ctx)
	}

	return password, string(bytes), nil
}

func (u *userSvcImpl) checkUserPassword(ctx context.Context, password, storedHash string) error {
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return errors.ErrUserInvalidPassword(ctx)
	}
	return nil
}

func (u *userSvcImpl) Create(ctx context.Context, user *auth.User) (*auth.User, error) {
	l := u.l().C(ctx).Mth("create").Trc()

	// check email
	if user.Username == "" {
		return nil, errors.ErrUserEmailEmpty(ctx)
	}
	if !kit.IsEmailValid(user.Username) {
		return nil, errors.ErrUserNoValidEmail(ctx)
	}

	user.Id = kit.NewId()
	user.LockedAt = nil

	// currently create activated users
	now := kit.Now()
	user.ActivatedAt = &now

	// check username uniqueness
	another, err := u.storage.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if another != nil {
		return nil, errors.ErrUserNameNotUnique(ctx, user.Username)
	}

	// check and hash password
	_, user.Password, err = u.getPwdHash(ctx, user.Password)
	if err != nil {
		return nil, err
	}

	// save to storage
	err = u.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	l.F(log.FF{"id": user.Id}).Dbg("created")

	return user, nil
}

func (u *userSvcImpl) GetByEmail(ctx context.Context, email string) (*auth.User, error) {
	u.l().C(ctx).Mth("get-by-username").F(log.FF{"email": email}).Trc()
	return u.storage.GetByUsername(ctx, email)
}

func (u *userSvcImpl) Get(ctx context.Context, userId string) (*auth.User, error) {
	u.l().C(ctx).Mth("get").F(log.FF{"userId": userId}).Trc()
	return u.storage.GetUser(ctx, userId)
}

func (u *userSvcImpl) GetByIds(ctx context.Context, userIds []string) ([]*auth.User, error) {
	u.l().C(ctx).Mth("get-ids").Trc()
	return u.storage.GetUserByIds(ctx, userIds)
}

func (u *userSvcImpl) SetPassword(ctx context.Context, userId, prevPassword, newPassword string) error {
	l := u.l().C(ctx).Mth("reset-password").F(log.FF{"userId": userId}).Trc()

	// find user
	user, err := u.storage.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound(ctx, userId)
	}

	//check user status
	if user.ActivatedAt == nil {
		return errors.ErrUserNotActive(ctx, userId)
	}
	if user.LockedAt != nil {
		return errors.ErrUserLocked(ctx, userId)
	}

	// check current password
	err = u.checkUserPassword(ctx, prevPassword, user.Password)
	if err != nil {
		return err
	}

	// get new password hash
	_, hash, err := u.getPwdHash(ctx, newPassword)
	if err != nil {
		return err
	}
	user.Password = hash

	if err := u.storage.UpdateUser(ctx, user); err != nil {
		return err
	}
	l.Trc("updated")
	return nil
}
