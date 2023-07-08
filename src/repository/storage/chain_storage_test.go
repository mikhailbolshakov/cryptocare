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

type chainStorageTestSuite struct {
	kitTestSuite.Suite
	storage domain.ChainStorage
	adapter Adapter
}

func (s *chainStorageTestSuite) SetupSuite() {
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

func (s *chainStorageTestSuite) TearDownSuite() {
	_ = s.adapter.Close(s.Ctx)
}

func TestChainStorageSuite(t *testing.T) {
	suite.Run(t, new(chainStorageTestSuite))
}

func (s *chainStorageTestSuite) getBid() *domain.Bid {
	currencies := []string{"RUB", "USD", "EUR", "BTC", "ETH", "SLN", "UAH", "CHY", "USDT", "USDC"}
	exchanges := []string{"binance", "huobi", "bybit"}
	iSrc := rand.Int31n(int32(len(currencies)))
	iTrg := rand.Int31n(int32(len(currencies)))
	iExch := rand.Int31n(int32(len(exchanges)))
	return &domain.Bid{
		Id:           kit.NewId(),
		Type:         domain.BidTypeP2P,
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

func (s *chainStorageTestSuite) getChain() *domain.ProfitableChain {
	bids := []*domain.Bid{s.getBid(), s.getBid()}
	return &domain.ProfitableChain{
		Id:            kit.NewRandString(),
		Asset:         kit.NewRandString(),
		ProfitShare:   1.12,
		Methods:       []string{"M1", "M2"},
		BidAssets:     []string{bids[0].SrcAsset, bids[0].TrgAsset, bids[1].TrgAsset},
		Bids:          bids,
		Depth:         2,
		ExchangeCodes: []string{bids[0].ExchangeCode, bids[1].ExchangeCode},
		CreatedAt:     kit.Now(),
	}
}

func (s *chainStorageTestSuite) Test_PutGet() {

	chains := []*domain.ProfitableChain{
		s.getChain(),
		s.getChain(),
	}

	err := s.storage.SaveProfitableChains(s.Ctx, chains)
	if err != nil {
		s.Fatal(err)
	}

	// get all
	rs, err := s.storage.GetProfitableChains(s.Ctx, &domain.GetProfitableChainsRequest{
		PagingRequest: kit.PagingRequest{Size: 10},
		Assets:        []string{chains[0].Asset, chains[1].Asset},
		WithBids:      true,
	})
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(rs.Chains)
	s.Equal(2, len(rs.Chains))
	for _, ch := range rs.Chains {
		s.NotEmpty(ch.Id)
		s.NotEmpty(ch.Asset)
		s.NotEmpty(ch.ProfitShare)
		s.NotEmpty(ch.CreatedAt)
		s.NotEmpty(ch.Bids)
		s.NotEmpty(ch.Methods)
		s.NotEmpty(ch.ExchangeCodes)
		s.NotEmpty(ch.BidAssets)
		s.NotEmpty(ch.Depth)
	}

	// get all
	rs, err = s.storage.GetProfitableChains(s.Ctx, &domain.GetProfitableChainsRequest{
		PagingRequest: kit.PagingRequest{Size: 10},
		Assets:        []string{chains[0].Asset},
		WithBids:      false,
	})
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(rs.Chains)
	s.Equal(1, len(rs.Chains))
	for _, ch := range rs.Chains {
		s.NotEmpty(ch.Id)
		s.Nil(ch.Bids)
	}

	// get single
	chain, err := s.storage.GetProfitableChain(s.Ctx, chains[0].Id)
	if err != nil {
		s.Fatal(err)
	}
	s.NotEmpty(chain)
	s.Equal(chain.Id, chains[0].Id)
	s.Equal(chain.Depth, chains[0].Depth)
	s.Equal(chain.ExchangeCodes, chains[0].ExchangeCodes)
	s.Equal(chain.BidAssets, chains[0].BidAssets)
	s.Equal(chain.ProfitShare, chains[0].ProfitShare)
	s.Equal(chain.Asset, chains[0].Asset)
	s.Equal(len(chain.Bids), len(chains[0].Bids))

	// unknown key
	chain, err = s.storage.GetProfitableChain(s.Ctx, kit.NewRandString())
	if err != nil {
		s.Fatal(err)
	}
	s.Nil(chain)

}
