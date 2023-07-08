//go:build integration
// +build integration

package storage

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type subscriptionStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.SubscriptionStorage
	adapter Adapter
}

func (s *subscriptionStorageTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())

	// load config
	cfg, err := service.LoadConfig()
	if err != nil {
		s.Fatal(err)
	}

	// initialize adapter
	s.adapter = NewAdapter()
	err = s.adapter.Init(s.Ctx, cfg)
	if err != nil {
		s.Fatal(err)
	}
	s.storage = s.adapter
	rand.Seed(time.Now().UnixNano())
}

func (s *subscriptionStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestSubscriptionStorageSuite(t *testing.T) {
	suite.Run(t, new(subscriptionStorageTestSuite))
}

func (s *subscriptionStorageTestSuite) getSubscription() *domain.Subscription {
	return &domain.Subscription{
		Id:       kit.NewId(),
		UserId:   kit.NewRandString(),
		IsActive: true,
		Filter: &domain.SubscriptionChainFilter{
			Assets:    []string{"RUB", "USD"},
			Methods:   []string{"M1", "M2"},
			Exchanges: []string{"binance", "huobi"},
			MaxDepth:  5,
			MinProfit: 0.5,
		},
		Notifications: []*domain.SubscriptionNotification{
			{
				Id:       kit.NewRandString(),
				Channel:  domain.SubscriptionNotificationChannelTelegram,
				IsActive: true,
				Telegram: &domain.SubscriptionTelegramNotificationDetails{
					Channel: 99999,
				},
			},
		},
	}
}

func (s *subscriptionStorageTestSuite) Test_CRUD() {

	subs := s.getSubscription()
	// create new
	err := s.storage.SaveSubscription(s.Ctx, subs)
	if err != nil {
		s.Fatal(err)
	}
	// get
	actual, err := s.storage.GetSubscription(s.Ctx, subs.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(actual)
	s.Equal(actual, subs)
	// update
	subs.Filter.MinProfit = 0.7
	subs.Filter.MaxDepth = 7
	subs.Filter.Assets = append(subs.Filter.Assets, "EUR")
	err = s.storage.SaveSubscription(s.Ctx, subs)
	if err != nil {
		s.Fatal(err)
	}
	// get
	actual, err = s.storage.GetSubscription(s.Ctx, subs.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(actual)
	s.Equal(actual, subs)
	// search active
	searchRs, err := s.storage.SearchSubscriptions(s.Ctx, &domain.SearchSubscriptionsRequest{})
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(searchRs)
	found := false
	for _, ss := range searchRs {
		if ss.Id == subs.Id {
			found = true
			break
		}
	}
	s.True(found)
	// delete
	err = s.storage.DeleteSubscription(s.Ctx, subs.Id)
	if err != nil {
		s.Fatal(err)
	}
	// get after delete
	actual, err = s.storage.GetSubscription(s.Ctx, subs.Id)
	if err != nil {
		s.Fatal(err)
	}
	s.Empty(actual)
}
