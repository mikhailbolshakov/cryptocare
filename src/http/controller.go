package http

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/context"
	kitHttp "github.com/mikhailbolshakov/cryptocare/src/kit/http"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"net/http"
	"strings"
)

type Controller interface {

	// Ready returns OK if service is ready
	Ready(http.ResponseWriter, *http.Request)

	// auth
	Login(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	Registration(http.ResponseWriter, *http.Request)
	TokenRefresh(http.ResponseWriter, *http.Request)
	SetPassword(http.ResponseWriter, *http.Request)

	// GetProfitableChains retrieves stored profit chains
	GetProfitableChains(http.ResponseWriter, *http.Request)
	// GetProfitableChainDetails retrieves details of the chain
	GetProfitableChainDetails(http.ResponseWriter, *http.Request)

	// subscriptions
	CreateSubscription(http.ResponseWriter, *http.Request)
	UpdateSubscription(http.ResponseWriter, *http.Request)
	GetSubscription(http.ResponseWriter, *http.Request)
	DeleteSubscription(http.ResponseWriter, *http.Request)
	GetUserSubscriptions(http.ResponseWriter, *http.Request)

	// bids
	PutBid(http.ResponseWriter, *http.Request)
}

type controllerIml struct {
	kitHttp.BaseController
	arbitrageService    domain.ArbitrageService
	userService         domain.UserService
	sessionService      auth.SessionsService
	subscriptionService domain.SubscriptionService
	bidProvider         domain.BidProvider
}

func NewController(arbitrageService domain.ArbitrageService, sessionService auth.SessionsService,
	userService domain.UserService, subscriptionService domain.SubscriptionService, bidProvider domain.BidProvider) Controller {
	return &controllerIml{
		BaseController: kitHttp.BaseController{
			Logger: service.LF(),
		},
		arbitrageService:    arbitrageService,
		sessionService:      sessionService,
		userService:         userService,
		subscriptionService: subscriptionService,
		bidProvider:         bidProvider,
	}
}

func (c *controllerIml) l() log.CLogger {
	return service.L().Cmp("controller")
}

// Ready godoc
// @Summary check system is ready
// @Router /ready [get]
// @Success 200
// @tags system
func (c *controllerIml) Ready(w http.ResponseWriter, r *http.Request) {
	c.RespondWithStatus(w, http.StatusOK, "OK")
}

// GetProfitableChains godoc
// @Summary retrieves profitable deal chains by criteria
// @Accept json
// @Produce json
// @Router /arbitrage/chains [get]
// @Param assets query string false "comma separated list of assets"
// @Param withBids query bool false "if chains are retrieved with bid info"
// @Param size query int false "page size"
// @Success 200 {object} ProfitableChains
// @Failure 500 {object} http.Error
// @tags arbitrage
func (c *controllerIml) GetProfitableChains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("get-chains").Trc()

	rq := &domain.GetProfitableChainsRequest{}

	assetsStr, err := c.FormVal(r, ctx, "assets", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if assetsStr != "" {
		rq.Assets = strings.Split(assetsStr, ",")
	}

	size, err := c.FormValInt(r, ctx, "size", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if size != nil {
		rq.PagingRequest.Size = *size
	}

	withBids, err := c.FormValBool(r, ctx, "withBids", true)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if withBids != nil {
		rq.WithBids = *withBids
	}

	chainsRs, err := c.arbitrageService.GetProfitableChains(ctx, rq)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	c.RespondOK(w, c.toProfitableChainsApi(chainsRs.Chains))
}

// GetProfitableChainDetails godoc
// @Summary retrieves profitable deal chain details by id
// @Accept json
// @Produce json
// @Router /arbitrage/chains/{chainId}/details [get]
// @Param chainId path string true "chain id"
// @Success 200 {object} ProfitableChain
// @Failure 500 {object} http.Error
// @tags arbitrage
func (c *controllerIml) GetProfitableChainDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.l().C(ctx).Mth("get-chain-details").Trc()

	chainId, err := c.Var(r, ctx, "chainId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	chain, err := c.arbitrageService.GetProfitableChain(ctx, chainId)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	c.RespondOK(w, c.toProfitableChainApi(chain))
}

// Registration godoc
// @Summary registers a new client
// @Accept json
// @produce json
// @Param regRequest body ClientRegistrationRequest true "registration request"
// @Success 200 {object} ClientUser
// @Failure 500 {object} http.Error
// @Router /auth/registration [post]
// @tags auth
func (c *controllerIml) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq := &ClientRegistrationRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	user, err := c.userService.Create(ctx, c.toClientRegRequestDomain(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toClientUserApi(user))
}

// Login godoc
// @Summary logins user by email/password
// @Accept json
// @produce json
// @Param loginRequest body LoginRequest true "auth request"
// @Success 200 {object} LoginResponse
// @Failure 500 {object} http.Error
// @Router /auth/login [post]
// @tags auth
func (c *controllerIml) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq := &LoginRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	session, token, err := c.sessionService.LoginPassword(ctx, c.toLoginRequest(rq))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toLoginResponseApi(session, token))
}

// Logout godoc
// @Summary logouts user
// @Accept json
// @Produce json
// @Router /auth/logout [post]
// @Success 200
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// take userId from the currently logged session
	rq, err := context.MustRequest(ctx)
	if err != nil || rq.Sid == "" {
		c.RespondError(w, errors.ErrLogoutNoSID(ctx))
		return
	}

	err = c.sessionService.Logout(ctx, rq.Sid)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// TokenRefresh godoc
// @Summary refreshes auth token
// @Accept json
// @Produce json
// @Router /auth/token/refresh [post]
// @Success 200 {object} SessionToken
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) TokenRefresh(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	_, err := context.MustRequest(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		c.RespondError(w, errors.ErrNoAuthHeader(ctx))
		return
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) < 2 {
		c.RespondError(w, errors.ErrAuthHeaderInvalid(ctx))
		return
	}
	token := splitToken[1]

	sessionToken, err := c.sessionService.RefreshToken(ctx, token)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toTokenApi(sessionToken))
}

// SetPassword godoc
// @Summary sets a new password for the user
// @Param request body SetPasswordRequest true "set password request"
// @Accept json
// @Produce json
// @Router /auth/password [post]
// @Success 200
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @tags auth
func (c *controllerIml) SetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// take userId from the currently logged session
	rqCtx, err := context.MustRequest(ctx)
	if err != nil {
		c.RespondError(w, err)
		return
	}
	if rqCtx.Uid == "" {
		c.RespondError(w, errors.ErrNoUID(ctx))
		return
	}

	rq := &SetPasswordRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	err = c.userService.SetPassword(ctx, rqCtx.Uid, rq.PrevPassword, rq.NewPassword)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, kitHttp.EmptyOkResponse)
}

// CreateSubscription godoc
// @Summary creates a new user subscription
// @Accept json
// @produce json
// @Param userId path string true "user id"
// @Param regRequest body SubscriptionRequest true "subscription request"
// @Success 200 {object} Subscription
// @Failure 500 {object} http.Error
// @Router /users/{userId}/subscriptions [post]
// @tags subscription
func (c *controllerIml) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := c.VarUUID(r, ctx, "userId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	if appCtx, ok := context.Request(ctx); ok && appCtx.GetUserId() != userId {
		c.RespondError(w, errors.ErrNotAllowed(ctx))
		return
	}

	rq := &SubscriptionRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	subscription, err := c.subscriptionService.Create(ctx, c.toCreateSubscriptionRequestDomain(rq, userId))
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toSubscriptionApi(subscription))
}

// UpdateSubscription godoc
// @Summary updates a subscription
// @Accept json
// @produce json
// @Param userId path string true "user id"
// @Param subscriptionId path string true "user id"
// @Param regRequest body SubscriptionRequest true "subscription request"
// @Success 200 {object} Subscription
// @Failure 500 {object} http.Error
// @Router /users/{userId}/subscriptions/{subscriptionId} [put]
// @tags subscription
func (c *controllerIml) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := c.VarUUID(r, ctx, "userId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	if appCtx, ok := context.Request(ctx); ok && appCtx.GetUserId() != userId {
		c.RespondError(w, errors.ErrNotAllowed(ctx))
		return
	}

	subsId, err := c.VarUUID(r, ctx, "subscriptionId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	rq := &SubscriptionRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	updRq := c.toCreateSubscriptionRequestDomain(rq, userId)
	updRq.Id = subsId
	subscription, err := c.subscriptionService.Update(ctx, updRq)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toSubscriptionApi(subscription))
}

func (c *controllerIml) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

// GetUserSubscriptions godoc
// @Summary retrieves user's subscriptions
// @Accept json
// @produce json
// @Param userId path string true "user id"
// @Success 200 {object} Subscriptions
// @Failure 500 {object} http.Error
// @Router /users/{userId}/subscriptions [get]
// @tags subscription
func (c *controllerIml) GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := c.VarUUID(r, ctx, "userId", false)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	if appCtx, ok := context.Request(ctx); ok && appCtx.GetUserId() != userId {
		c.RespondError(w, errors.ErrNotAllowed(ctx))
		return
	}

	subscriptions, err := c.subscriptionService.Search(ctx, &domain.SearchSubscriptionsRequest{UserId: userId})
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toSubscriptionsApi(subscriptions))
}

func (c *controllerIml) GetSubscription(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

// PutBid godoc
// @Summary allows creation or updating an exchange bid
// @Accept json
// @produce json
// @Param request body BidRequest true "bid request"
// @Success 200 {object} Bid
// @Failure 500 {object} http.Error
// @Router /arbitrage/bids [post]
// @tags subscription
func (c *controllerIml) PutBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rq := &BidRequest{}
	if err := c.DecodeRequest(r, ctx, rq); err != nil {
		c.RespondError(w, err)
		return
	}

	bidRq := c.toBidRequestDomain(rq)
	bid, err := c.bidProvider.PutBid(ctx, bidRq)
	if err != nil {
		c.RespondError(w, err)
		return
	}

	c.RespondOK(w, c.toBidApi(bid))
}
