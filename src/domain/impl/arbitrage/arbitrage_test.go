package arbitrage

import (
	_ "embed"
	"encoding/json"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/mocks"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type arbitrageTestSuite struct {
	kitTestSuite.Suite
	bidsProvider *mocks.BidProvider
	chainStorage *mocks.ChainStorage
	notifier     *mocks.Notifier
	svc          domain.ArbitrageService
}

func (s *arbitrageTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func TestArbitrageSuite(t *testing.T) {
	suite.Run(t, new(arbitrageTestSuite))
}

func (s *arbitrageTestSuite) SetupTest() {
	s.bidsProvider = &mocks.BidProvider{}
	s.chainStorage = &mocks.ChainStorage{}
	s.notifier = &mocks.Notifier{}
	s.svc = NewArbitrageService(s.chainStorage, s.bidsProvider, s.notifier)
	s.svc.Init(&service.Config{Arbitrage: &service.Arbitrage{Depth: 5, MinProfit: 1.0005, CheckLimit: true}})
}

var (
	//go:embed arbitrage_test_find_chains_data.json
	orderChainsTestData []byte
)

func (s *arbitrageTestSuite) ChainsToStr(c *domain.CandidateChains) []string {
	r := []string{}
	for _, chain := range c.Chains {
		bidStr := ""
		for _, bidId := range chain.BidIds {
			bidStr += bidId + "->"
		}
		if bidStr != "" {
			r = append(r, bidStr)
		}
	}
	return r
}

func (s *arbitrageTestSuite) Test_FindChains() {

	svc := s.svc.(*arbitrageSvcImpl)

	var tests []*struct {
		Name     string             `json:"name"`
		Bids     []*domain.BidLight `json:"bids"`
		Asset    string             `json:"asset"`
		Expected []string           `json:"expectedChains"`
	}
	_ = json.Unmarshal(orderChainsTestData, &tests)

	for _, tt := range tests {
		s.T().Run(tt.Name, func(t *testing.T) {
			bidsMap := make(map[string][]*domain.BidLight)
			for _, b := range tt.Bids {
				bidsMap[b.SrcAsset] = append(bidsMap[b.SrcAsset], b)
			}
			s.bidsProvider.ExpectedCalls = nil
			m := s.bidsProvider.On("GetBidLightsBySourceAsset", s.Ctx, mock.AnythingOfType("string"))
			m.RunFn = func(args mock.Arguments) {
				m.ReturnArguments = mock.Arguments{bidsMap[args.Get(1).(string)], nil}
			}
			actual := &domain.CandidateChains{}
			err := svc.findChainsRecurse(s.Ctx, tt.Asset, tt.Asset, nil, actual, 0)
			s.Nil(err)
			actualStr := s.ChainsToStr(actual)
			s.Equal(actualStr, tt.Expected)
		})
	}

}

func (s *arbitrageTestSuite) Test_BuildProfitableChains_WhenEmptyCandidates_Empty_Ok() {
	svc := s.svc.(*arbitrageSvcImpl)
	var candidates []*domain.CandidateChain
	profitableChains, err := svc.buildProfitableChains(s.Ctx, candidates)
	s.Nil(err)
	s.Empty(profitableChains)
}

func (s *arbitrageTestSuite) Test_BuildProfitableChains_WhenSingleNewChain_Ok() {
	svc := s.svc.(*arbitrageSvcImpl)
	candidates := []*domain.CandidateChain{
		{
			BidIds:    []string{kit.NewRandString(), kit.NewRandString()},
			TotalRate: 1.1,
		},
	}
	bids := []*domain.Bid{
		{
			Id:           candidates[0].BidIds[0],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "USD",
			TrgAsset:     "RUB",
			Rate:         63,
			ExchangeCode: "binance",
			Methods:      []string{"M1", "M2"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[0].BidIds[1],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "RUB",
			TrgAsset:     "USD",
			Rate:         0.015,
			ExchangeCode: "bitnami",
			Methods:      []string{"M2", "M3"},
			UserId:       kit.NewId(),
			Link:         "",
		},
	}
	s.bidsProvider.On("GetBidsByIds", s.Ctx, candidates[0].BidIds).Return(bids, nil)
	s.chainStorage.On("ProfitableChainExists", s.Ctx, mock.AnythingOfType("string")).Return(false, nil)
	profitableChains, err := svc.buildProfitableChains(s.Ctx, candidates)
	s.Nil(err)
	s.Len(profitableChains, 1)
	s.NotEmpty(profitableChains[0].Id)
	s.ElementsMatch([]string{"M1", "M2", "M3"}, profitableChains[0].Methods)
	s.ElementsMatch([]string{"binance", "bitnami"}, profitableChains[0].ExchangeCodes)
	s.Equal([]string{"USD", "RUB", "USD"}, profitableChains[0].BidAssets)
	s.Equal("USD", profitableChains[0].Asset)
	s.Equal(2, profitableChains[0].Depth)
	s.Equal(candidates[0].TotalRate, profitableChains[0].ProfitShare)
	s.NotEmpty(profitableChains[0].CreatedAt)
	s.Len(profitableChains[0].Bids, 2)
}

func (s *arbitrageTestSuite) Test_BuildProfitableChains_WhenChainExists_Ok() {
	svc := s.svc.(*arbitrageSvcImpl)
	candidates := []*domain.CandidateChain{
		{
			BidIds:    []string{kit.NewRandString(), kit.NewRandString()},
			TotalRate: 1.1,
		},
	}
	bids := []*domain.Bid{
		{
			Id:           candidates[0].BidIds[0],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "USD",
			TrgAsset:     "RUB",
			Rate:         63,
			ExchangeCode: "binance",
			Methods:      []string{"M1", "M2"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[0].BidIds[1],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "RUB",
			TrgAsset:     "USD",
			Rate:         0.015,
			ExchangeCode: "bitnami",
			Methods:      []string{"M2", "M3"},
			UserId:       kit.NewId(),
			Link:         "",
		},
	}
	s.bidsProvider.On("GetBidsByIds", s.Ctx, candidates[0].BidIds).Return(bids, nil)
	s.chainStorage.On("ProfitableChainExists", s.Ctx, mock.AnythingOfType("string")).Return(true, nil)
	profitableChains, err := svc.buildProfitableChains(s.Ctx, candidates)
	s.Nil(err)
	s.Empty(profitableChains)
}

func (s *arbitrageTestSuite) Test_BuildProfitableChains_WhenDuplicatedNewChains_Ok() {
	svc := s.svc.(*arbitrageSvcImpl)
	candidates := []*domain.CandidateChain{
		{
			BidIds:    []string{kit.NewRandString(), kit.NewRandString()},
			TotalRate: 1.1,
		},
	}
	candidates = append(candidates, candidates[0])
	bids := []*domain.Bid{
		{
			Id:           candidates[0].BidIds[0],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "USD",
			TrgAsset:     "RUB",
			Rate:         63,
			ExchangeCode: "binance",
			Methods:      []string{"M1", "M2"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[0].BidIds[1],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "RUB",
			TrgAsset:     "USD",
			Rate:         0.015,
			ExchangeCode: "bitnami",
			Methods:      []string{"M2", "M3"},
			UserId:       kit.NewId(),
			Link:         "",
		},
	}
	s.bidsProvider.On("GetBidsByIds", s.Ctx, candidates[0].BidIds).Return(bids, nil)
	s.bidsProvider.On("GetBidsByIds", s.Ctx, candidates[1].BidIds).Return(bids, nil)
	s.chainStorage.On("ProfitableChainExists", s.Ctx, mock.AnythingOfType("string")).Return(false, nil)
	profitableChains, err := svc.buildProfitableChains(s.Ctx, candidates)
	s.Nil(err)
	s.Len(profitableChains, 1)
}

func (s *arbitrageTestSuite) Test_BuildProfitableChains_WhenMultipleNewChains_Ok() {
	svc := s.svc.(*arbitrageSvcImpl)
	candidates := []*domain.CandidateChain{
		{
			BidIds:    []string{kit.NewRandString(), kit.NewRandString()},
			TotalRate: 1.1,
		},
		{
			BidIds:    []string{kit.NewRandString(), kit.NewRandString()},
			TotalRate: 1.1,
		},
	}
	bids := []*domain.Bid{
		{
			Id:           candidates[0].BidIds[0],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "USD",
			TrgAsset:     "RUB",
			Rate:         63,
			ExchangeCode: "binance",
			Methods:      []string{"M1", "M2"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[0].BidIds[1],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "RUB",
			TrgAsset:     "USD",
			Rate:         0.015,
			ExchangeCode: "bitnami",
			Methods:      []string{"M2", "M3"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[1].BidIds[0],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "USD",
			TrgAsset:     "RUB",
			Rate:         63,
			ExchangeCode: "binance",
			Methods:      []string{"M1", "M2"},
			UserId:       kit.NewId(),
			Link:         "",
		},
		{
			Id:           candidates[1].BidIds[1],
			Type:         domain.BidTypeP2P,
			SrcAsset:     "RUB",
			TrgAsset:     "USD",
			Rate:         0.015,
			ExchangeCode: "bitnami",
			Methods:      []string{"M2", "M3"},
			UserId:       kit.NewId(),
			Link:         "",
		},
	}
	s.bidsProvider.On("GetBidsByIds", s.Ctx, append(candidates[0].BidIds, candidates[1].BidIds...)).Return(bids, nil)
	s.chainStorage.On("ProfitableChainExists", s.Ctx, mock.AnythingOfType("string")).Return(false, nil)
	profitableChains, err := svc.buildProfitableChains(s.Ctx, candidates)
	s.Nil(err)
	s.Len(profitableChains, 2)
}
