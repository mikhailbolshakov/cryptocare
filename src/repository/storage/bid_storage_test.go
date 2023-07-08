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

type bidStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.BidStorage
	adapter Adapter
}

func (s *bidStorageTestSuite) SetupSuite() {
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

func (s *bidStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestBidStorageSuite(t *testing.T) {
	suite.Run(t, new(bidStorageTestSuite))
}

func (s *bidStorageTestSuite) getBid() *domain.Bid {
	currencies := []string{"RUB", "USD", "EUR", "BTC", "ETH", "SLN", "UAH", "CHY", "USDT", "USDC"}
	exchanges := []string{"binance", "huobi", "bybit"}
	iSrc := rand.Int31n(int32(len(currencies)))
	iTrg := rand.Int31n(int32(len(currencies)))
	iExch := rand.Int31n(int32(len(exchanges)))
	return &domain.Bid{
		Id:           kit.NewId(),
		SrcAsset:     currencies[iSrc],
		TrgAsset:     currencies[iTrg],
		Rate:         rand.Float64(),
		ExchangeCode: exchanges[iExch],
		Available:    rand.Float64(),
		MinLimit:     10.5,
		MaxLimit:     1000.5,
		Methods:      []string{"M1", "M2", "M3"},
		Link:         "http://link",
		UserId:       kit.NewRandString(),
	}
}

func (s *bidStorageTestSuite) Test_PutGet() {

	// generate multiple bids
	bids := make([]*domain.Bid, 100)
	for i := 0; i < 100; i++ {
		bids[i] = s.getBid()
	}
	// put to store
	err := s.storage.PutBids(s.Ctx, bids)
	if err != nil {
		s.Fatal(err)
	}

	// get all
	bidsLight, err := s.storage.GetBidsLightAll(s.Ctx)
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(bidsLight)
	ids := make([]string, len(bidsLight))
	for i, b := range bidsLight {
		s.NotEmpty(b.Id)
		s.NotEmpty(b.SrcAsset)
		s.NotEmpty(b.TrgAsset)
		s.NotEmpty(b.Rate)
		s.NotEmpty(b.Available)
		s.NotEmpty(b.MaxLimit)
		s.NotEmpty(b.MinLimit)
		ids[i] = b.Id
	}

	// get by ids
	bids, err = s.storage.GetBidsByIds(s.Ctx, ids)
	if err != nil {
		s.Fatal(err)
	}
	s.Equal(len(bids), len(ids))
}
