//go:build integration
// +build integration

package storage

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type sessionStorageTestSuite struct {
	kitTestSuite.Suite
	storage auth.SessionStorage
	adapter Adapter
}

func (s *sessionStorageTestSuite) SetupSuite() {

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

func (s *sessionStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func (s *sessionStorageTestSuite) getSession() *auth.Session {
	id := kit.NewId()
	now := kit.Now()
	return &auth.Session{
		Id:             id,
		UserId:         kit.NewId(),
		Username:       id + "@test.com",
		LoginAt:        now,
		LogoutAt:       nil,
		LastActivityAt: now,
		Roles:          []string{"client"},
	}
}

func (s *sessionStorageTestSuite) SetupTest() {}

func TestSessionSuite(t *testing.T) {
	suite.Run(t, new(sessionStorageTestSuite))
}

func (s *sessionStorageTestSuite) Test_Session_CRUD() {

	// create sessions
	expected1 := s.getSession()
	err := s.storage.CreateSession(s.Ctx, expected1)
	if err != nil {
		s.Fatal(err)
	}

	expected2 := s.getSession()
	expected2.UserId = expected1.UserId
	err = s.storage.CreateSession(s.Ctx, expected2)
	if err != nil {
		s.Fatal(err)
	}

	// get session
	actual, err := s.storage.Get(s.Ctx, expected1.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected1, actual)

	// get session by id from cache
	actual, err = s.storage.Get(s.Ctx, expected1.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected1, actual)

	// get sessions by user
	actuals, err := s.storage.GetByUser(s.Ctx, expected1.UserId)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(2, len(actuals))

	now := kit.Now()
	// update session
	expected1.LastActivityAt = now
	err = s.storage.UpdateLastActivity(s.Ctx, expected1.Id, expected1.LastActivityAt)
	if err != nil {
		s.Fatal(err)
	}

	// logout
	expected1.LogoutAt = &now
	err = s.storage.Logout(s.Ctx, expected1.Id, now)
	if err != nil {
		s.Fatal(err)
	}

	// get session
	actual, err = s.storage.Get(s.Ctx, expected1.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected1, actual)

	// get session by id from cache
	actual, err = s.storage.Get(s.Ctx, expected1.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(expected1, actual)

	// get sessions (only active)
	actuals, err = s.storage.GetByUser(s.Ctx, expected1.UserId)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(1, len(actuals))

}
