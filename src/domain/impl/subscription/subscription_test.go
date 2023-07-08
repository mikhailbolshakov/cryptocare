package subscription

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/mocks"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type subscriptionTestSuite struct {
	kitTestSuite.Suite
	storage  *mocks.SubscriptionStorage
	notifier *mocks.TelegramNotifier
	svc      domain.SubscriptionService
}

func (s *subscriptionTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func TestSubscriptionTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(subscriptionTestSuite))
}

func (s *subscriptionTestSuite) SetupTest() {
	s.storage = &mocks.SubscriptionStorage{}
	s.notifier = &mocks.TelegramNotifier{}
	s.svc = NewSubscriptionService(s.storage, s.notifier)
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Depth: 5, MinProfit: 1.0005}})
}

func (s *subscriptionTestSuite) getSubscription() *domain.Subscription {
	return &domain.Subscription{
		Id:       kit.NewId(),
		UserId:   kit.NewId(),
		IsActive: true,
		Filter: &domain.SubscriptionChainFilter{
			Assets:    []string{"USD", "RUB"},
			Methods:   []string{"M1", "M2"},
			Exchanges: []string{"binance", "bybit"},
			MaxDepth:  5,
			MinProfit: 1,
		},
		Notifications: []*domain.SubscriptionNotification{
			{
				Id:       kit.NewRandString(),
				Channel:  domain.SubscriptionNotificationChannelTelegram,
				IsActive: true,
				Telegram: &domain.SubscriptionTelegramNotificationDetails{
					Channel: -124125123515,
				},
			},
		},
	}
}

func (s *subscriptionTestSuite) Test_ValidateAndPopulate_WhenProcess_Ok() {
	subs := s.getSubscription()
	subs.Filter.Assets = []string{" Rub ", " USD "}
	subs.Filter.Exchanges = []string{" biNance  ", " bybit   "}
	err := s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.Nil(err)
	s.Equal(subs.Filter.Assets, []string{"RUB", "USD"})
	s.Equal(subs.Filter.Exchanges, []string{"binance", "bybit"})
}

func (s *subscriptionTestSuite) Test_ValidateAndPopulate_WhenMinProfitInvalid_Fail() {
	subs := s.getSubscription()
	subs.Filter.MinProfit = 0.000000001
	err := s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionMinProfitInvalid)
	subs.Filter.MinProfit = -1.0
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionMinProfitInvalid)
	subs.Filter.MinProfit = 101.0
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionMinProfitInvalid)
	subs.Filter.MinProfit = 10.0
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.Nil(err)
}

func (s *subscriptionTestSuite) Test_ValidateAndPopulate_WhenMaxDepthInvalid_Fail() {
	subs := s.getSubscription()
	subs.Filter.MaxDepth = 1
	err := s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionMaxDepthInvalid)
	subs.Filter.MaxDepth = -1
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionMaxDepthInvalid)
	subs.Filter.MaxDepth = 0
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.Nil(err)
}

func (s *subscriptionTestSuite) Test_ValidateAndPopulate_WhenNotificationInvalid_Fail() {
	subs := s.getSubscription()
	subs.Notifications[0].Channel = "unknown"
	err := s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionNotificationChannelNotSupported)
	subs = s.getSubscription()
	subs.Notifications[0].Telegram.Channel = 0
	err = s.svc.(*subscriptionSvcImpl).validateAndPopulate(s.Ctx, subs)
	s.AssertAppErr(err, errors.ErrCodeSubscriptionNotificationTelegramInvalid)
}

func (s *subscriptionTestSuite) Test_Create_Ok() {
	subs := s.getSubscription()
	var actual *domain.Subscription
	s.storage.On("SaveSubscription", s.Ctx, subs).
		Run(func(args mock.Arguments) {
			actual = args.Get(1).(*domain.Subscription)
		}).Return(nil)
	_, err := s.svc.Create(s.Ctx, subs)
	s.Nil(err)
	s.NotEmpty(actual)
	s.NotEmpty(actual.Id)
	s.True(actual.IsActive)
}

func (s *subscriptionTestSuite) Test_Update_Ok() {
	stored := s.getSubscription()
	subs := s.getSubscription()
	subs.Id = stored.Id
	subs.UserId = stored.UserId
	subs.Filter.Exchanges = []string{"another"}
	s.storage.On("GetSubscription", s.Ctx, stored.Id).Return(stored, nil)
	var actual *domain.Subscription
	s.storage.On("SaveSubscription", s.Ctx, subs).
		Run(func(args mock.Arguments) {
			actual = args.Get(1).(*domain.Subscription)
		}).Return(nil)
	_, err := s.svc.Update(s.Ctx, subs)
	s.Nil(err)
	s.NotEmpty(actual)
	s.Equal(actual.Id, stored.Id)
	s.Equal(actual.UserId, stored.UserId)
	s.True(actual.IsActive)
	s.Equal(actual.Filter.Exchanges, []string{"another"})
}

func (s *subscriptionTestSuite) Test_Notify_OneChainOneSubscriptionMatch_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "RUB",
			ProfitShare:   1.2,
			Methods:       []string{"M1", "M2"},
			Depth:         2,
			ExchangeCodes: []string{"exch1"},
			CreatedAt:     time.Time{},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1", "exch2"}
	sub1.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub1.Filter.Methods = []string{"M1", "M2", "M3"}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 1
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Equal(len(actualChannels), 1)
	s.Equal(len(actualChains), 1)
}

func (s *subscriptionTestSuite) Test_Notify_OneChainOneSubscriptionDoesntMatchByAsset_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "UAH",
			ProfitShare:   1.2,
			Methods:       []string{"M1", "M2"},
			Depth:         2,
			ExchangeCodes: []string{"exch1"},
			CreatedAt:     time.Time{},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1", "exch2"}
	sub1.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub1.Filter.Methods = []string{"M1", "M2", "M3"}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 1
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Empty(actualChannels)
	s.Empty(actualChains)
}

func (s *subscriptionTestSuite) Test_Notify_OneChainOneSubscriptionDoesntMatchByMethods_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "UAH",
			ProfitShare:   1.2,
			Methods:       []string{"M1", "M2", "M4"},
			Depth:         2,
			ExchangeCodes: []string{"exch1"},
			CreatedAt:     time.Time{},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1", "exch2"}
	sub1.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub1.Filter.Methods = []string{"M1", "M2", "M3"}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 1
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Empty(actualChannels)
	s.Empty(actualChains)
}

func (s *subscriptionTestSuite) Test_Notify_OneChainOneSubscriptionDoesntMatchByMinProfit_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "UAH",
			ProfitShare:   1.08,
			Methods:       []string{"M1", "M2", "M4"},
			Depth:         2,
			ExchangeCodes: []string{"exch1"},
			CreatedAt:     time.Time{},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1", "exch2"}
	sub1.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub1.Filter.Methods = []string{"M1", "M2", "M3"}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 10
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Empty(actualChannels)
	s.Empty(actualChains)
}

func (s *subscriptionTestSuite) Test_Notify_TwoChainTwoSubscriptionMatch_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "RUB",
			ProfitShare:   1.2,
			Methods:       []string{"M1", "M2"},
			Depth:         2,
			ExchangeCodes: []string{"exch1", "exch2"},
		},
		{
			Id:            kit.NewId(),
			Asset:         "EUR",
			ProfitShare:   1.2,
			Methods:       []string{"M3", "M4"},
			Depth:         2,
			ExchangeCodes: []string{"exch2", "exch3"},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1", "exch2", "exch3"}
	sub1.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub1.Filter.Methods = []string{"M1", "M2", "M3", "M4"}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 10
	sub2 := s.getSubscription()
	sub2.Filter.Exchanges = []string{"exch1", "exch2", "exch3"}
	sub2.Filter.Assets = []string{"RUB", "USD", "EUR"}
	sub2.Filter.Methods = []string{"M1", "M2", "M3", "M4"}
	sub2.Filter.MaxDepth = 5
	sub2.Filter.MinProfit = 10
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1, sub2}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Equal(len(actualChannels), 4)
	s.Equal(len(actualChains), 2)
}

func (s *subscriptionTestSuite) Test_Notify_OneChainOneSubscription_Match_MethodsSanitized_Ok() {
	chains := []*domain.ProfitableChain{
		{
			Id:            kit.NewId(),
			Asset:         "RUB",
			ProfitShare:   1.2,
			Methods:       []string{"  M_  %*^%*^ == 1 ++  "},
			Depth:         2,
			ExchangeCodes: []string{"exch1"},
		},
	}
	sub1 := s.getSubscription()
	sub1.Filter.Exchanges = []string{"exch1"}
	sub1.Filter.Assets = []string{"RUB"}
	sub1.Filter.Methods = []string{"----    M^(&^(*&^( 1....  "}
	sub1.Filter.MaxDepth = 5
	sub1.Filter.MinProfit = 10
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Notification: &service.ArbitrageNotification{Telegram: &service.ArbitrageNotificationTelegram{Bot: "bot"}}}})
	var actualChannels []int
	var actualChains []*domain.ProfitableChain
	s.notifier.On("Notify", s.Ctx, mock.AnythingOfType("string"), mock.AnythingOfType("[]int"), mock.AnythingOfType("[]*domain.ProfitableChain")).
		Run(func(args mock.Arguments) {
			actualChannels = append(actualChannels, args.Get(2).([]int)...)
			actualChains = append(actualChains, args.Get(3).([]*domain.ProfitableChain)...)
		}).
		Return(nil)
	s.storage.On("SearchSubscriptions", s.Ctx, mock.AnythingOfType("*domain.SearchSubscriptionsRequest")).Return([]*domain.Subscription{sub1}, nil)
	err := s.svc.Notify(s.Ctx, chains)
	s.Nil(err)
	s.Equal(len(actualChannels), 1)
	s.Equal(len(actualChains), 1)
}
