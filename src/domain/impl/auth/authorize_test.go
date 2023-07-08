package auth

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/mocks"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type authorizeTestSuite struct {
	kitTestSuite.Suite
	storage      *mocks.SessionStorage
	authorizeSvc auth.AuthorizeSession
}

func (s *authorizeTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func (s *authorizeTestSuite) SetupTest() {
	s.storage = &mocks.SessionStorage{}
	s.authorizeSvc = NewAuthorizeService(s.storage)
}

func TestAuthorizeSuite(t *testing.T) {
	suite.Run(t, new(authorizeTestSuite))
}

func (s *authorizeTestSuite) Test_Authorize_When_FoundInCache_Ok() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "resource",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	key := s.authorizeSvc.(*authorizeSvcImpl).authorizeSessionCacheKey(rq)
	// preset cache: allowed
	s.authorizeSvc.(*authorizeSvcImpl).cache.Set(key, true, time.Hour)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.True(allowed)
	// preset cache: not allowed
	s.authorizeSvc.(*authorizeSvcImpl).cache.Set(key, false, time.Hour)
	allowed, err = s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.False(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_SessionNotFound_Fail() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "resource",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(nil, nil)
	_, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.AssertAppErr(err, errors.ErrCodeSessionNotFound)
}

func (s *authorizeTestSuite) Test_Authorize_When_SessionLogOut_Fail() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "resource",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	now := kit.Now()
	session := &auth.Session{LogoutAt: &now}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	_, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.AssertAppErr(err, errors.ErrCodeSessionLoggedOut)
}

func (s *authorizeTestSuite) Test_Authorize_When_NoSessionRoles_PermissionDenied() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "resource",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	session := &auth.Session{}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.False(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_SysadminRole() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "resource",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleSysAdmin},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.True(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_ClientNoPermissionConfigured() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    "unknown",
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleArbitrageClient},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	_, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.AssertAppErr(err, errors.ErrCodeSessionAuthorizationInvalidResource)
}

func (s *authorizeTestSuite) Test_Authorize_When_ClientNoPermission_AccessDenied() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    domain.AuthResUserProfileAll,
				Permissions: []string{auth.AccessR, auth.AccessW},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleArbitrageClient},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.False(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_ClientGrantedPermission() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    domain.AuthResUserProfileMy,
				Permissions: []string{auth.AccessR},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleArbitrageClient},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.True(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_ClientNoPermission_MultipleAccess() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    domain.AuthResUserProfileMy,
				Permissions: []string{auth.AccessR, auth.AccessD},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleArbitrageClient},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.False(allowed)
}

func (s *authorizeTestSuite) Test_Authorize_When_ClientNoPermission_MultipleResources() {
	rq := &auth.AuthorizationRequest{
		SessionId: kit.NewRandString(),
		AuthorizationResources: []*auth.AuthorizationResource{
			{
				Resource:    domain.AuthResUserProfileMy,
				Permissions: []string{auth.AccessR},
			},
			{
				Resource:    domain.AuthResUserProfileMy,
				Permissions: []string{auth.AccessD},
			},
		},
	}
	session := &auth.Session{
		Roles: []string{domain.AuthRoleArbitrageClient},
	}
	s.storage.On("Get", s.Ctx, rq.SessionId).Return(session, nil)
	allowed, err := s.authorizeSvc.AuthorizeSession(s.Ctx, rq)
	s.Nil(err)
	s.False(allowed)
}

func (s *authorizeTestSuite) Test_GetRolesForGroups() {
	// two groups
	roles, err := s.authorizeSvc.GetRolesForGroups(s.Ctx, []string{domain.AuthGroupSysAdmin, domain.AuthGroupClient})
	s.Nil(err)
	s.ElementsMatch(roles, []string{domain.AuthRoleSysAdmin, domain.AuthRoleArbitrageClient})
	// one group
	roles, err = s.authorizeSvc.GetRolesForGroups(s.Ctx, []string{domain.AuthGroupClient})
	s.Nil(err)
	s.ElementsMatch(roles, []string{domain.AuthRoleArbitrageClient})
	// empty
	roles, err = s.authorizeSvc.GetRolesForGroups(s.Ctx, []string{})
	s.Nil(err)
	s.ElementsMatch(roles, []string{})

}
