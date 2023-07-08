package auth

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	memcache "github.com/mikhailbolshakov/cryptocare/src/kit/cache"
	kitError "github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"strings"
	"time"
)

type authorizeSvcImpl struct {
	storage auth.SessionStorage
	cache   memcache.MemCache
}

func NewAuthorizeService(storage auth.SessionStorage) auth.AuthorizeSession {
	return &authorizeSvcImpl{
		storage: storage,
		cache:   memcache.NewMemCache(),
	}
}

func (s *authorizeSvcImpl) l() log.CLogger {
	return service.L().Cmp("authorize-svc")
}

func (s *authorizeSvcImpl) authorizeSessionCacheKey(rq *auth.AuthorizationRequest) string {
	sb := strings.Builder{}
	sb.WriteString(rq.SessionId)
	for _, r := range rq.AuthorizationResources {
		sb.WriteString(r.Resource)
		sb.WriteString(strings.Join(r.Permissions, ""))
	}
	return sb.String()
}

type rolePermissions struct {
	Role        string
	Permissions []string
}

// permissions specifies access on resources for session roles
var permissions = map[string][]rolePermissions{
	domain.AuthResUserProfileAll:     {rolePermissions{Role: domain.AuthRoleSysAdmin, Permissions: []string{auth.AccessR, auth.AccessW, auth.AccessD}}},
	domain.AuthResArbitrageChainsAll: {rolePermissions{Role: domain.AuthRoleArbitrageClient, Permissions: []string{auth.AccessR}}},
	domain.AuthResUserProfileMy:      {rolePermissions{Role: domain.AuthRoleArbitrageClient, Permissions: []string{auth.AccessR, auth.AccessW}}},
}

func (s *authorizeSvcImpl) authorizeSession(ctx context.Context, rq *auth.AuthorizationRequest) error {
	s.l().C(ctx).Mth("authorize-int").Trc()

	// get session by sid
	session, err := s.storage.Get(ctx, rq.SessionId)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.ErrSessionNotFound(ctx)
	}

	// check session is logged out
	if session.LogoutAt != nil {
		return errors.ErrSessionLoggedOut(ctx)
	}

	// check permissions
	if len(session.Roles) == 0 {
		return errors.ErrSecurityPermissionsDenied(ctx)
	}

	// don't check permissions for sysadmin
	if kit.Strings(session.Roles).Contains(domain.AuthRoleSysAdmin) {
		return nil
	}

	// check for all requested resources
	for _, rqResource := range rq.AuthorizationResources {
		// check resource permissions
		resourcePerms, ok := permissions[rqResource.Resource]
		// resource not found
		if !ok {
			return errors.ErrSessionAuthorizationInvalidResource(ctx)
		}
		granted := false
		// check if there is a role for session that is granted with the requested permissions
		for _, rolePerm := range resourcePerms {
			if kit.Strings(session.Roles).Contains(rolePerm.Role) &&
				kit.Strings(rqResource.Permissions).Intersect(rolePerm.Permissions).Equal(rqResource.Permissions) {
				granted = true
				break
			}
		}
		if !granted {
			return errors.ErrSecurityPermissionsDenied(ctx)
		}
	}
	return nil
}

func (s *authorizeSvcImpl) AuthorizeSession(ctx context.Context, rq *auth.AuthorizationRequest) (bool, error) {
	l := s.l().C(ctx).Mth("authorize").Trc()

	if rq.SessionId == "" {
		return false, errors.ErrSidEmpty(ctx)
	}

	// note, permissions for session stay the same along the session life
	key := s.authorizeSessionCacheKey(rq)

	// get from cache
	if v, ok := s.cache.Get(key); ok {
		l.Trc("found in cache")
		return v.(bool), nil
	} else {
		// if no cache hit, check authorization for session
		if err := s.authorizeSession(ctx, rq); err != nil {
			if appErr, ok := kitError.Is(err); ok && appErr.Code() == errors.ErrCodeSecurityPermissionsDenied {
				s.cache.Set(key, false, time.Hour)
				return false, nil
			}
			return false, err
		} else {
			s.cache.Set(key, true, time.Hour)
			return true, nil
		}
	}
}

var groupRoles = map[string][]string{
	domain.AuthGroupSysAdmin: {domain.AuthRoleSysAdmin},
	domain.AuthGroupClient:   {domain.AuthRoleArbitrageClient},
}

func (s *authorizeSvcImpl) GetRolesForGroups(ctx context.Context, groups []string) ([]string, error) {
	s.l().C(ctx).Mth("roles-for-groups").F(log.FF{"groups": groups}).Trc()
	var r []string
	for _, gr := range groups {
		if roles, ok := groupRoles[gr]; ok {
			r = append(r, roles...)
		}
	}
	return kit.Strings(r).Distinct(), nil
}
