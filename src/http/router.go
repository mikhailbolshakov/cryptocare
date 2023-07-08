package http

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth/impl"
	"github.com/mikhailbolshakov/cryptocare/src/kit/http"
	_ "github.com/mikhailbolshakov/cryptocare/src/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	ctrl         Controller
	routeBuilder *http.RouteBuilder
}

func NewRouter(c Controller, routeBuilder *http.RouteBuilder) http.RouteSetter {
	return &Router{
		ctrl:         c,
		routeBuilder: routeBuilder,
	}
}

func (r *Router) Set() error {
	return r.routeBuilder.Build(
		// readiness
		http.R("/ready", r.ctrl.Ready).GET().NoAuth(),

		// authentication & authorization
		http.R("/api/auth/login", r.ctrl.Login).POST().NoAuth(),
		http.R("/api/auth/token/refresh", r.ctrl.TokenRefresh).POST().NoAuth(),
		http.R("/api/auth/logout", r.ctrl.Logout).POST(),
		http.R("/api/auth/registration", r.ctrl.Registration).POST().NoAuth(),
		http.R("/api/auth/password", r.ctrl.SetPassword).POST(),

		// subscriptions
		http.R("/api/users/{userId}/subscriptions", r.ctrl.CreateSubscription).POST(),
		http.R("/api/users/{userId}/subscriptions/{subscriptionId}", r.ctrl.UpdateSubscription).PUT(),
		http.R("/api/users/{userId}/subscriptions/{subscriptionId}", r.ctrl.DeleteSubscription).DELETE(),
		http.R("/api/users/{userId}/subscriptions/{subscriptionId}", r.ctrl.GetSubscription).GET(),
		http.R("/api/users/{userId}/subscriptions", r.ctrl.GetUserSubscriptions).GET(),

		// arbitrage
		http.R("/api/arbitrage/chains", r.ctrl.GetProfitableChains).GET().Authorize(impl.Resource(domain.AuthResArbitrageChainsAll, "r")),
		http.R("/api/arbitrage/chains/{chainId}/details", r.ctrl.GetProfitableChainDetails).GET().Authorize(impl.Resource(domain.AuthResArbitrageChainsAll, "r")),

		// bids
		http.R("/api/arbitrage/bids", r.ctrl.PutBid).POST(),

		// swagger
		http.R("", nil).PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler),
	)
}
