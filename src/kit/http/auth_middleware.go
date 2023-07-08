package http

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"net/http"
	"strings"
)

type Middleware struct {
	BaseController
	authSessionRepository      auth.AuthenticateSession
	authorizeSessionRepository auth.AuthorizeSession
	resourcePolicyManager      auth.ResourcePolicyManager
}

func NewMiddleware(logger log.CLoggerFunc, authSessionRepository auth.AuthenticateSession,
	authorizeSessionRepository auth.AuthorizeSession, resourcePolicyManager auth.ResourcePolicyManager) *Middleware {
	return &Middleware{
		authSessionRepository:      authSessionRepository,
		authorizeSessionRepository: authorizeSessionRepository,
		resourcePolicyManager:      resourcePolicyManager,
		BaseController: BaseController{
			Logger: logger,
		},
	}
}

func (m *Middleware) AuthAccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {

		// check if context is a request context
		ctxRq, err := context.MustRequest(r.Context())
		if err != nil {
			m.RespondError(w, err)
			return
		}
		ctx := r.Context()

		// check and extract Authorization data
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.RespondError(w, ErrSecurityLoginFailed(ctx))
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) < 2 {
			m.RespondError(w, ErrSecurityLoginFailed(ctx))
			return
		}
		accessToken := splitToken[1]

		// authorize session
		session, err := m.authSessionRepository.AuthSession(ctx, accessToken)
		if err != nil {
			m.RespondError(w, ErrSecurityLoginFailed(ctx))
			return
		}

		if session == nil {
			m.RespondError(w, ErrSecurityLoginFailed(ctx))
			return
		}

		// populate context based on login params
		ctx = ctxRq.
			WithUser(session.UserId, session.Username).
			WithSessionId(session.Id).
			ToContext(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return f
}

func (m *Middleware) SetContextMiddleware(next http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		// init context
		ctxRq := context.NewRequestCtx().Rest()

		// check and set Request ID coming from client
		requestId := r.Header.Get("RequestId")
		if requestId != "" {
			ctxRq = ctxRq.WithRequestId(requestId)
		} else {
			ctxRq = ctxRq.WithNewRequestId()
		}

		ctx := ctxRq.ToContext(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (m *Middleware) AuthorizationMiddleware(routeId string, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		if rqCtx, ok := context.Request(ctx); ok {
			sid := rqCtx.GetSessionId()

			// build authorization request for the resources and permissions
			authorizationResources, err := m.resourcePolicyManager.GetRequestedResources(ctx, routeId, r)
			if err != nil {
				m.RespondError(w, err)
				return
			}

			// if no resources requested, error
			if len(authorizationResources) == 0 {
				m.RespondError(w, ErrSecurityPermissionsDenied(ctx))
				return
			}

			// authorize session
			allowed, err := m.authorizeSessionRepository.AuthorizeSession(ctx, &auth.AuthorizationRequest{
				SessionId:              sid,
				AuthorizationResources: authorizationResources,
			})
			if err != nil {
				m.RespondError(w, err)
				return
			}
			if !allowed {
				m.RespondError(w, ErrSecurityPermissionsDenied(ctx))
				return
			}

		}

		next(w, r)
	}
}
