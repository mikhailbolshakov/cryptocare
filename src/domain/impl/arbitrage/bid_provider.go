package arbitrage

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/goroutine"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"go.uber.org/atomic"
	"strings"
	"sync"
	"time"
)

type bidProviderImpl struct {
	sync.RWMutex
	bidStorage        domain.BidStorage
	bidLightsMap      map[string][]*domain.BidLight
	assets            map[string]struct{}
	assetsRestriction map[string]struct{}
	cancelFunc        context.CancelFunc
	running           *atomic.Bool
	cfg               *service.Config
}

func NewBidProviderService(bidStorage domain.BidStorage) domain.BidProvider {
	return &bidProviderImpl{
		bidStorage:        bidStorage,
		running:           atomic.NewBool(false),
		assetsRestriction: make(map[string]struct{}),
	}
}

func (s *bidProviderImpl) l() log.CLogger {
	return service.L().Cmp("bid-provider")
}

func (s *bidProviderImpl) Init(cfg *service.Config) {
	s.cfg = cfg
	if s.cfg.Arbitrage.Assets != "" {
		for _, a := range strings.Split(s.cfg.Arbitrage.Assets, ",") {
			s.assetsRestriction[a] = struct{}{}
		}
	}
}

func (s *bidProviderImpl) Run(ctx context.Context) error {
	l := s.l().C(ctx).Mth("run").Trc()

	// check running
	if s.running.Load() {
		return errors.ErrBidProviderAlreadyRun(ctx)
	}

	ctx, s.cancelFunc = context.WithCancel(ctx)
	s.running.Store(true)

	goroutine.New().
		WithLogger(s.l().C(ctx).Mth("chains-notify-worker")).
		WithRetry(goroutine.Unrestricted).
		WithRetryDelay(time.Second*10).
		Go(ctx, func() {
			ticker := time.NewTicker(time.Duration(s.cfg.Arbitrage.BidProviderPeriodSec) * time.Second)
			lg := s.l().C(ctx).Mth("get-bids")
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					bids, err := s.bidStorage.GetBidsLightAll(ctx)
					if err != nil {
						s.l().C(ctx).Mth("get-bids").E(err).Err()
						continue
					}
					lg.DbgF("found: %d", len(bids))
					bidLights := make(map[string][]*domain.BidLight)
					assets := make(map[string]struct{})
					for _, b := range bids {
						bidLights[b.SrcAsset] = append(bidLights[b.SrcAsset], b)
						// populate assets
						if len(s.assetsRestriction) == 0 {
							assets[b.SrcAsset] = struct{}{}
						} else {
							if _, ok := s.assetsRestriction[b.SrcAsset]; ok {
								assets[b.SrcAsset] = struct{}{}
							}
						}
					}

					// swap newly read data with stored
					s.Lock()
					s.bidLightsMap = bidLights
					s.assets = assets
					s.Unlock()

				case <-ctx.Done():
					l.Inf("stop")
					return
				}
			}
		})
	return nil
}

func (s *bidProviderImpl) Stop(ctx context.Context) error {
	l := s.l().C(ctx).Mth("stop").Trc()
	// cancel if running
	if s.cancelFunc != nil && s.running.Load() {
		s.cancelFunc()
		s.running.Store(false)
		s.cancelFunc = nil
		l.Inf("ok")
	}
	return nil
}

func (s *bidProviderImpl) GetAssets(ctx context.Context) ([]string, error) {
	l := s.l().C(ctx).Mth("get-assets")
	s.RLock()
	defer s.RUnlock()
	var r []string
	for k := range s.assets {
		r = append(r, k)
	}
	l.TrcF("found: %d", len(r))
	return r, nil
}

func (s *bidProviderImpl) GetBidLightsBySourceAsset(ctx context.Context, srcAsset string) ([]*domain.BidLight, error) {
	s.RLock()
	defer s.RUnlock()
	return s.bidLightsMap[srcAsset], nil
}

func (s *bidProviderImpl) GetBidsByIds(ctx context.Context, ids []string) ([]*domain.Bid, error) {
	return s.bidStorage.GetBidsByIds(ctx, ids)
}

func (s *bidProviderImpl) PutBid(ctx context.Context, bid *domain.Bid) (*domain.Bid, error) {
	s.l().C(ctx).Mth("put").Trc()
	if bid.Id == "" {
		bid.Id = kit.NewRandString()
	}
	bid.Type = domain.BidTypeManual
	err := s.bidStorage.PutBids(ctx, []*domain.Bid{bid}, 60*60*4)
	if err != nil {
		return nil, err
	}
	return bid, nil
}
