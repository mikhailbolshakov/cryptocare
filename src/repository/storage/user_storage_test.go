//go:build integration
// +build integration

package storage

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type userStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.UserStorage
	adapter Adapter
}

func (s *userStorageTestSuite) SetupSuite() {

	s.Suite.Init(service.LF())

	// load config
	cfg, err := service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	// disable applying migrations
	cfg.Storages.Pg.MigPath = ""

	// initialize adapter
	s.adapter = NewAdapter()
	err = s.adapter.Init(s.Ctx, cfg)
	if err != nil {
		s.Fatal(err)
	}
	s.storage = s.adapter
}

func (s *userStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func (s *userStorageTestSuite) getUser() *auth.User {
	id := kit.NewId()
	now := kit.Now()
	return &auth.User{
		Id:          id,
		Username:    id + "@test.com",
		Password:    kit.NewRandString(),
		Type:        domain.UserTypeClient,
		FirstName:   "First",
		LastName:    "Last",
		ActivatedAt: &now,
		Groups:      []string{domain.AuthGroupClient},
	}
}

func (s *userStorageTestSuite) SetupTest() {}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userStorageTestSuite))
}

func (s *userStorageTestSuite) Test_User_CRUD() {

	// create a user
	expected := s.getUser()
	err := s.storage.CreateUser(s.Ctx, expected)
	if err != nil {
		s.Fatal(err)
	}

	// get user by id
	actual, err := s.storage.GetUser(s.Ctx, expected.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// get user by id from cache
	actual, err = s.storage.GetUser(s.Ctx, expected.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// get user by username
	actual, err = s.storage.GetByUsername(s.Ctx, expected.Username)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// update user (set locked)
	now := kit.Now()
	expected.LockedAt = &now
	err = s.storage.UpdateUser(s.Ctx, expected)
	if err != nil {
		s.Fatal(err)
	}

	// get updated
	actual, err = s.storage.GetUser(s.Ctx, expected.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// get user by id from cache
	actual, err = s.storage.GetUser(s.Ctx, expected.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// get user by username
	actual, err = s.storage.GetByUsername(s.Ctx, expected.Username)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected, actual)

	// get by ids
	actuals, err := s.storage.GetUserByIds(s.Ctx, []string{expected.Id})
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(1, len(actuals))
	s.Equal(expected, actuals[0])

	// delete user
	err = s.storage.DeleteUser(s.Ctx, actual)
	if err != nil {
		s.Fatal(err)
	}

	// get user by id from cache
	actual, err = s.storage.GetUser(s.Ctx, expected.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Nil(actual)

	actual, err = s.storage.GetByUsername(s.Ctx, expected.Username)
	if err != nil {
		s.Fatal(err)
	}
	s.Nil(actual)

	actuals, err = s.storage.GetUserByIds(s.Ctx, []string{expected.Id})
	if err != nil {
		s.Fatal(err)
	}
	s.Empty(actuals)

}
