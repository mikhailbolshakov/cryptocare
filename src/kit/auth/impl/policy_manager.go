package impl

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"net/http"
	"strings"
)

type resourcePolicyManager struct {
	routePoliciesMap map[string][]auth.ResourcePolicy
	logger           log.CLoggerFunc
}

type ConditionFn func(context.Context, *http.Request) (bool, error)

type ResourcePolicyBuilder struct {
	resource    string
	permissions []string
	conditions  []ConditionFn
}

func (s *resourcePolicyManager) l() log.CLogger {
	return s.logger().Cmp("resource-policy-manager")
}

func NewResourcePolicyManager(logger log.CLoggerFunc) auth.ResourcePolicyManager {
	return &resourcePolicyManager{
		routePoliciesMap: map[string][]auth.ResourcePolicy{},
		logger:           logger,
	}
}

func (s *resourcePolicyManager) RegisterResourceMapping(routeId string, policies ...auth.ResourcePolicy) {
	s.l().Mth("register").Trc(routeId, " registered")
	s.routePoliciesMap[routeId] = policies
}

func (s *resourcePolicyManager) GetRequestedResources(ctx context.Context, routeId string, r *http.Request) ([]*auth.AuthorizationResource, error) {
	l := s.l().Mth("get-resources")

	var resources []*auth.AuthorizationResource
	var codes []string

	if policies, ok := s.routePoliciesMap[routeId]; ok {
		for _, policy := range policies {
			resource, err := policy.Resolve(ctx, r)
			if err != nil {
				return nil, err
			}
			if resource == nil {
				continue
			}
			resources = append(resources, resource)
			codes = append(codes, resource.Resource)
		}
	}
	l.F(log.FF{"routeId": routeId, "resources": codes}).Trc()

	return resources, nil
}

func Resource(resource string, permissions string) *ResourcePolicyBuilder {
	b := &ResourcePolicyBuilder{
		resource:   resource,
		conditions: []ConditionFn{},
	}
	b.permissions = b.convertPermissions(permissions)
	return b
}

// convertPermissions converts permissions from "rwxd" string to []string{"r", w", "x", "d"}
func (a *ResourcePolicyBuilder) convertPermissions(permissions string) []string {
	var res []string
	s := strings.ToLower(permissions)
	if strings.Contains(s, auth.AccessR) {
		res = append(res, auth.AccessR)
	}
	if strings.Contains(s, auth.AccessW) {
		res = append(res, auth.AccessW)
	}
	if strings.Contains(s, auth.AccessX) {
		res = append(res, auth.AccessX)
	}
	if strings.Contains(s, auth.AccessD) {
		res = append(res, auth.AccessD)
	}
	return res
}

func (a *ResourcePolicyBuilder) When(f ConditionFn) *ResourcePolicyBuilder {
	a.conditions = append(a.conditions, f)
	return a
}

func (a *ResourcePolicyBuilder) WhenNot(f ConditionFn) *ResourcePolicyBuilder {
	a.conditions = append(a.conditions, func(c context.Context, r *http.Request) (bool, error) { res, err := f(c, r); return !res, err })
	return a
}

func (a *ResourcePolicyBuilder) Resolve(ctx context.Context, r *http.Request) (*auth.AuthorizationResource, error) {
	// check conditions
	for _, cond := range a.conditions {
		if condRes, err := cond(ctx, r); err != nil {
			return nil, err
		} else {
			if !condRes {
				return nil, nil
			}
		}
	}
	return &auth.AuthorizationResource{
		Resource:    a.resource,
		Permissions: a.permissions,
	}, nil
}

func (a *ResourcePolicyBuilder) B() auth.ResourcePolicy {
	return a
}
