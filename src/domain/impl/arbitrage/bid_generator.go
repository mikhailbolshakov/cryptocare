package arbitrage

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/goroutine"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

type bidGeneratorImpl struct {
	bidStorage domain.BidStorage
	cancelFunc context.CancelFunc
	running    *atomic.Bool
	cfg        *service.Config
}

func NewBidGenerator(bidStorage domain.BidStorage) domain.BidGenerator {
	return &bidGeneratorImpl{
		bidStorage: bidStorage,
		running:    atomic.NewBool(false),
	}
}

func (b *bidGeneratorImpl) l() log.CLogger {
	return service.L().Cmp("bid-gen")
}

func (b *bidGeneratorImpl) Init(cfg *service.Config) {
	b.cfg = cfg
	rand.Seed(time.Now().UnixNano())
}

var (
	currencies  = []string{"RUB", "USD", "EUR", "BTC", "ETH", "SLN", "USDT", "AVL"}
	exchanges   = []string{"binance", "huobi", "bybit"}
	rateRandCfg = map[string]struct {
		Min float64
		Max float64
	}{
		"RUB-USD":  {Min: 0.015, Max: 0.016},
		"RUB-EUR":  {Min: 0.014, Max: 0.015},
		"RUB-BTC":  {Min: 0.00000075, Max: 0.00000078},
		"RUB-ETH":  {Min: 0.000010600, Max: 0.000010650},
		"RUB-SLN":  {Min: 0.0005, Max: 0.0006},
		"RUB-USDT": {Min: 0.015, Max: 0.016},
		"RUB-AVL":  {Min: 0.0008, Max: 0.0009},

		"USD-RUB":  {Min: 54.45, Max: 56.5},
		"USD-EUR":  {Min: 1.017, Max: 1.019},
		"USD-BTC":  {Min: 0.000048, Max: 0.000049},
		"USD-ETH":  {Min: 0.00065, Max: 0.00067},
		"USD-SLN":  {Min: 0.031, Max: 0.033},
		"USD-USDT": {Min: 0.95, Max: 0.98},
		"USD-AVL":  {Min: 0.053, Max: 0.055},

		"EUR-RUB":  {Min: 57.8, Max: 58.0},
		"EUR-USD":  {Min: 0.95, Max: 0.97},
		"EUR-BTC":  {Min: 0.000047, Max: 0.000048},
		"EUR-ETH":  {Min: 0.00064, Max: 0.00065},
		"EUR-SLN":  {Min: 0.031, Max: 0.032},
		"EUR-USDT": {Min: 0.9, Max: 0.95},
		"EUR-AVL":  {Min: 0.05, Max: 0.06},

		"BTC-RUB":  {Min: 1280000, Max: 1280050},
		"BTC-USD":  {Min: 20300, Max: 20330},
		"BTC-EUR":  {Min: 20010, Max: 20015},
		"BTC-ETH":  {Min: 13.67057, Max: 13.67058},
		"BTC-SLN":  {Min: 645.165, Max: 645.168},
		"BTC-USDT": {Min: 19800, Max: 19810},
		"BTC-AVL":  {Min: 1095.955, Max: 1095.96},

		"ETH-RUB":  {Min: 90957.34, Max: 90957.80},
		"ETH-USD":  {Min: 1488.7, Max: 1488.8},
		"ETH-EUR":  {Min: 1451, Max: 1451.5},
		"ETH-BTC":  {Min: 0.073, Max: 0.074},
		"ETH-SLN":  {Min: 47.15700, Max: 47.15750},
		"ETH-USDT": {Min: 1488.7, Max: 1488.8},
		"ETH-AVL":  {Min: 80.150, Max: 80.160},

		"SLN-RUB":  {Min: 1920, Max: 1930},
		"SLN-USD":  {Min: 31.3, Max: 31.4},
		"SLN-EUR":  {Min: 31, Max: 31.2},
		"SLN-BTC":  {Min: 0.00154, Max: 0.00157},
		"SLN-ETH":  {Min: 0.021, Max: 0.023},
		"SLN-USDT": {Min: 31.3, Max: 31.5},
		"SLN-AVL":  {Min: 1.69, Max: 1.71},

		"USDT-RUB": {Min: 54.45, Max: 58.5},
		"USDT-EUR": {Min: 1.017, Max: 1.019},
		"USDT-BTC": {Min: 0.000048, Max: 0.000050},
		"USDT-ETH": {Min: 0.00065, Max: 0.00070},
		"USDT-SLN": {Min: 0.031, Max: 0.033},
		"USDT-USD": {Min: 0.95, Max: 1.1},
		"USDT-AVL": {Min: 0.0536, Max: 0.06},

		"AVL-RUB":  {Min: 1100, Max: 1200},
		"AVL-USD":  {Min: 18.21, Max: 18.3},
		"AVL-EUR":  {Min: 19, Max: 19.5},
		"AVL-BTC":  {Min: 0.00085, Max: 0.0009},
		"AVL-ETH":  {Min: 0.0124, Max: 0.0128},
		"AVL-SLN":  {Min: 0.588, Max: 0.595},
		"AVL-USDT": {Min: 18.21, Max: 18.35},
	}
)

func (b *bidGeneratorImpl) getBid() *domain.Bid {
	iSrc := rand.Int31n(int32(len(currencies)))
	iTrg := rand.Int31n(int32(len(currencies)))
	iExch := rand.Int31n(int32(len(exchanges)))
	if iSrc == iTrg {
		return b.getBid()
	}
	key := currencies[iSrc] + "-" + currencies[iTrg]
	min := rateRandCfg[key].Min
	max := rateRandCfg[key].Max
	rate := min + (max-min)*rand.Float64()
	minLimit := 10 * rand.Float64()
	maxLimit := 10 + (1000-10)*rand.Float64()
	availableMin := minLimit * rate
	availableMax := maxLimit * rate
	available := availableMin + (availableMax-availableMin)*rand.Float64()
	return &domain.Bid{
		Id:           kit.NewId(),
		Type:         domain.BidTypeP2P,
		SrcAsset:     currencies[iSrc],
		TrgAsset:     currencies[iTrg],
		Rate:         rate,
		ExchangeCode: exchanges[iExch],
		Available:    available,
		MinLimit:     minLimit,
		MaxLimit:     maxLimit,
		Methods:      []string{"M1", "M2", "M3"},
		Link:         fmt.Sprintf("https://binance.com/orders?order=%s", kit.NewRandString()),
		UserId:       kit.NewId(),
	}
}

func (b *bidGeneratorImpl) Run(ctx context.Context) {
	l := b.l().C(ctx).Mth("run").Trc()

	ctx, b.cancelFunc = context.WithCancel(ctx)
	b.running.Store(true)

	goroutine.New().
		WithLogger(l).
		WithRetry(goroutine.Unrestricted).
		WithRetryDelay(time.Second*10).
		Go(ctx, func() {
			ticker := time.NewTicker(time.Duration(b.cfg.Dev.BidGeneratorPeriodSec) * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// generate bids
					bids := make([]*domain.Bid, b.cfg.Dev.BidGeneratorBidsCount)
					for i := 0; i < b.cfg.Dev.BidGeneratorBidsCount; i++ {
						bids[i] = b.getBid()
					}
					err := b.bidStorage.PutBids(ctx, bids, 60*10)
					if err != nil {
						continue
					}
				case <-ctx.Done():
					l.Inf("stop")
					return
				}
			}
		})
}

func (b *bidGeneratorImpl) Stop(ctx context.Context) {
	l := b.l().C(ctx).Mth("stop").Trc()
	if b.cancelFunc != nil && b.running.Load() {
		b.cancelFunc()
		b.running.Store(false)
		b.cancelFunc = nil
		l.Inf("ok")
	}
}
