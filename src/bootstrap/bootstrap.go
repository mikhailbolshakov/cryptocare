package bootstrap

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/domain/impl/arbitrage"
	"github.com/mikhailbolshakov/cryptocare/src/domain/impl/auth"
	"github.com/mikhailbolshakov/cryptocare/src/domain/impl/subscription"
	"github.com/mikhailbolshakov/cryptocare/src/http"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth/impl"
	kitHttp "github.com/mikhailbolshakov/cryptocare/src/kit/http"
	kitService "github.com/mikhailbolshakov/cryptocare/src/kit/service"
	"github.com/mikhailbolshakov/cryptocare/src/kit/telegram"
	"github.com/mikhailbolshakov/cryptocare/src/repository/storage"
	"github.com/mikhailbolshakov/cryptocare/src/service"
)

type serviceImpl struct {
	cfg                 *service.Config
	http                *kitHttp.Server
	arbitrageService    domain.ArbitrageService
	bidProvider         domain.BidProvider
	storageAdapter      storage.Adapter
	bidTestGenerator    domain.BidGenerator
	subscriptionService domain.SubscriptionService
}

// New creates a new instance of the service
func New() kitService.Service {
	s := &serviceImpl{}

	s.storageAdapter = storage.NewAdapter()
	s.bidProvider = arbitrage.NewBidProviderService(s.storageAdapter)
	s.bidTestGenerator = arbitrage.NewBidGenerator(s.storageAdapter)

	return s
}

func (s *serviceImpl) GetCode() string {
	return "trading"
}

// Init does all initializations
func (s *serviceImpl) Init(ctx context.Context) error {

	// load config
	var err error
	s.cfg, err = service.LoadConfig()
	if err != nil {
		return err
	}

	// set log config
	service.Logger.Init(s.cfg.Log)

	telegramNotifier := subscription.NewTelegramNotifier(telegram.NewTelegram(service.LF()),
		&subscription.TelegramOptions{
			Bot: s.cfg.Arbitrage.Notification.Telegram.Bot,
		})
	s.subscriptionService = subscription.NewSubscriptionService(s.storageAdapter, telegramNotifier)
	s.arbitrageService = arbitrage.NewArbitrageService(s.storageAdapter, s.bidProvider, s.subscriptionService)

	// create HTTP server
	s.http = kitHttp.NewHttpServer(s.cfg.Http, service.LF())

	// create resource policy manager
	resourcePolicyManager := impl.NewResourcePolicyManager(service.LF())
	authorizeSession := auth.NewAuthorizeService(s.storageAdapter)
	sessionService := impl.NewSessionsService(service.LF(), s.storageAdapter, s.storageAdapter, authorizeSession)
	userService := auth.NewUserService(s.storageAdapter)

	// create and set middlewares
	mdw := kitHttp.NewMiddleware(service.LF(), sessionService, authorizeSession, resourcePolicyManager)
	s.http.RootRouter.Use(mdw.SetContextMiddleware)

	// set up routing
	routeBuilder := kitHttp.NewRouteBuilder(s.http, resourcePolicyManager, mdw)

	// setup routes & controllers
	routers := []kitHttp.RouteSetter{
		http.NewRouter(http.NewController(s.arbitrageService, sessionService, userService, s.subscriptionService, s.bidProvider), routeBuilder),
	}
	for _, r := range routers {
		if err := r.Set(); err != nil {
			return err
		}
	}

	// init services
	s.arbitrageService.Init(s.cfg)
	sessionService.Init(s.cfg.Auth)
	s.bidTestGenerator.Init(s.cfg)
	s.bidProvider.Init(s.cfg)
	s.subscriptionService.Init(s.cfg)
	_ = telegramNotifier.Init(ctx)

	if err := s.storageAdapter.Init(ctx, s.cfg); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Start(ctx context.Context) error {

	// start listening REST
	s.http.Listen()

	// run bids generator for development mode
	if s.cfg.Dev.Enabled {
		s.bidTestGenerator.Run(ctx)
	}

	// start background arbitrage
	if err := s.arbitrageService.RunCalculationBackground(ctx); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Close(ctx context.Context) {
	s.bidTestGenerator.Stop(ctx)
	_ = s.arbitrageService.StopCalculation(ctx)
	_ = s.storageAdapter.Close(ctx)
	s.http.Close()
}
