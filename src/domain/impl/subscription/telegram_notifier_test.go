package subscription

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/telegram"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/stretchr/testify/suite"
	"testing"
)

type telegramNotifierTestSuite struct {
	kitTestSuite.Suite
	svc domain.TelegramNotifier
}

func (s *telegramNotifierTestSuite) SetupSuite() {
	s.Suite.Init(service.LF())
}

func TestTelegramNotifierTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(telegramNotifierTestSuite))
}

func (s *telegramNotifierTestSuite) SetupTest() {
	s.svc = NewTelegramNotifier(telegram.NewTelegram(s.L), &TelegramOptions{
		Bot: "5751469157:AAGXGOcqaENcw7uL8HoI5wS_dEX_y6Hlk88",
	})
}

func (s *telegramNotifierTestSuite) Test() {
	svc := s.svc.(*telegramNotifier)
	chain := &domain.ProfitableChain{
		Id:          "2345325325235",
		Asset:       "USD",
		ProfitShare: 1.25,
		Methods:     []string{"M1", "M2"},
		BidAssets:   []string{"USD", "USDT"},
		Bids: []*domain.Bid{
			{
				Id:           "11",
				Type:         domain.BidTypeP2P,
				SrcAsset:     "USD",
				TrgAsset:     "USDT",
				Rate:         0.97,
				ExchangeCode: "binance",
				Methods:      []string{"M1", "M2"},
			},
		},
		ExchangeCodes: []string{"binance", "huobi"},
	}
	rq := svc.getRequest(chain)
	s.NotEmpty(rq)
	//err := svc.Notify(s.Ctx, []*domain.ProfitableChain{chain})
	//if err != nil {
	//	s.L().E(err).Err()
	//	s.Fatal(err)
	//}
}
