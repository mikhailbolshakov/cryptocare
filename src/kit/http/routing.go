package http

import (
	"github.com/gorilla/mux"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"net/http"
)

type RouteBuilder struct {
	resourcePolicyManager auth.ResourcePolicyManager
	http                  *Server
	mdw                   *Middleware
}

func NewRouteBuilder(http *Server, policyManager auth.ResourcePolicyManager, mdw *Middleware) *RouteBuilder {
	return &RouteBuilder{
		resourcePolicyManager: policyManager,
		http:                  http,
		mdw:                   mdw,
	}
}

type Route struct {
	id               string
	url              string
	urlPrefix        string
	verbs            []string
	handleFn         http.HandlerFunc
	handler          http.Handler
	auth             bool
	subRouter        bool
	resourcePolicies []auth.ResourcePolicy
	middlewares      []mux.MiddlewareFunc
}

func (r *RouteBuilder) Build(routes ...*Route) error {

	for _, route := range routes {
		// validate route
		err := route.validate()
		if err != nil {
			return err
		}

		// setup http routing
		httpRouter := r.http.RootRouter
		if route.subRouter {
			httpRouter = r.http.RootRouter.PathPrefix(route.urlPrefix).Handler(route.handler).Subrouter()
			// set middlewares is passed
			if len(route.middlewares) > 0 {
				httpRouter.Use(route.middlewares...)
			}
		}

		if route.handleFn != nil {
			handleFn := route.handleFn
			// if authentication, apply special middleware
			if route.auth {
				// if authorization resource policies specified, cover handle func with authorization handle func
				if len(route.resourcePolicies) > 0 {
					// register resource mapping
					r.resourcePolicyManager.RegisterResourceMapping(route.id, route.resourcePolicies...)
					// setup authorization middleware
					handleFn = r.mdw.AuthorizationMiddleware(route.id, route.handleFn)
				}
				// apply token auth middleware
				handleFn = r.mdw.AuthAccessTokenMiddleware(handleFn)
			}
			httpRouter.HandleFunc(route.url, handleFn).Methods(route.verbs...)
		} else if route.handler != nil {
			// if handler specified, it means all processing done by it
			httpRouter.PathPrefix(route.urlPrefix).Handler(route.handler)
		}
	}
	return nil
}

// R starts building a new route with url and handle function
func R(url string, f func(http.ResponseWriter, *http.Request)) *Route {
	return &Route{
		id:       kit.NewRandString(),
		url:      url,
		handleFn: f,
		auth:     true,
	}
}

// NoAuth marks route as not required authentication
func (r *Route) NoAuth() *Route {
	r.auth = false
	return r
}

// SubRouter allows specifying a new area of routes with its own set of middlewares
func (r *Route) SubRouter(urlPrefix string) *Route {
	r.subRouter = true
	r.urlPrefix = urlPrefix
	return r
}

// Middlewares allows specifying special middlewares applied to the route
// Note! It's applied only to SubRoute
func (r *Route) Middlewares(mdws ...mux.MiddlewareFunc) *Route {
	r.middlewares = append(r.middlewares, mdws...)
	return r
}

// Url specifies route's URL
func (r *Route) Url(url string) *Route {
	r.url = url
	return r
}

// PathPrefix specifies URL prefix
func (r *Route) PathPrefix(urlPrefix string) *Route {
	r.urlPrefix = urlPrefix
	return r
}

// POST applies post verb
func (r *Route) POST() *Route {
	r.verbs = append(r.verbs, "POST")
	return r
}

// PUT applies put verb
func (r *Route) PUT() *Route {
	r.verbs = append(r.verbs, "PUT")
	return r
}

// GET applies get verb
func (r *Route) GET() *Route {
	r.verbs = append(r.verbs, "GET")
	return r
}

// DELETE applies delete verb
func (r *Route) DELETE() *Route {
	r.verbs = append(r.verbs, "DELETE")
	return r
}

// HandleFn specifies a handle function for route
func (r *Route) HandleFn(f func(http.ResponseWriter, *http.Request)) *Route {
	r.handleFn = f
	return r
}

// Handler allows specifying a handler which is applied to URL prefix
func (r *Route) Handler(h http.Handler) *Route {
	r.handler = h
	return r
}

// Authorize allows specifying authorization policy
func (r *Route) Authorize(policies ...auth.ResourcePolicy) *Route {
	r.resourcePolicies = append(r.resourcePolicies, policies...)
	return r
}

func (r *Route) validate() error {
	if r.url == "" && r.urlPrefix == "" {
		return ErrRouteBuilderUrlEmpty()
	}
	if len(r.verbs) == 0 && r.handleFn != nil {
		return ErrRouteBuilderVerbEmpty(r.url)
	}
	if r.handler == nil && r.handleFn == nil {
		return ErrRouteBuilderBothHandleFuncAndHandlerEmpty(r.url)
	}
	if len(r.resourcePolicies) > 0 && !r.auth {
		return ErrRouteBuilderAuthorizationPoliciesSpecifiedWithoutAuthentication(r.url)
	}
	if len(r.middlewares) > 0 && !r.subRouter {
		return ErrRouteBuilderSpecialMiddlewaresRequireSubrouting(r.url)
	}
	return nil
}
