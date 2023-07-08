package arbitrage

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/goroutine"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"github.com/mitchellh/hashstructure/v2"
	"go.uber.org/atomic"
	"math"
	"strconv"
	"time"
)

type arbitrageSvcImpl struct {
	bidProvider                 domain.BidProvider
	chainStorage                domain.ChainStorage
	assetsToCalculateChan       chan string
	saveProfitableChainsChan    chan []*domain.ProfitableChain
	processProfitableChainsChan chan []*domain.CandidateChain
	profitableChainsNotifyChan  chan []*domain.ProfitableChain
	cancelFunc                  context.CancelFunc
	running                     *atomic.Bool
	cfg                         *service.Config
	notifier                    domain.Notifier
}

func NewArbitrageService(chainStorage domain.ChainStorage, bidProvider domain.BidProvider, notifier domain.Notifier) domain.ArbitrageService {
	return &arbitrageSvcImpl{
		chainStorage:                chainStorage,
		bidProvider:                 bidProvider,
		assetsToCalculateChan:       make(chan string, 10),
		processProfitableChainsChan: make(chan []*domain.CandidateChain, 10),
		saveProfitableChainsChan:    make(chan []*domain.ProfitableChain, 10),
		profitableChainsNotifyChan:  make(chan []*domain.ProfitableChain, 10),
		running:                     atomic.NewBool(false),
		notifier:                    notifier,
	}
}

func (s *arbitrageSvcImpl) l() log.CLogger {
	return service.L().Cmp("arbitrage-svc")
}

func (s *arbitrageSvcImpl) Init(cfg *service.Config) {
	s.cfg = cfg
}

func (s *arbitrageSvcImpl) copyChain(chain *domain.CandidateChain) *domain.CandidateChain {
	r := &domain.CandidateChain{
		BidIds:    make([]string, len(chain.BidIds)),
		TotalRate: chain.TotalRate,
		Amount:    chain.Amount,
	}
	copy(r.BidIds, chain.BidIds)
	return r
}

// findChainsRecurse is a recursive func used for calculating one stage of deals
func (s *arbitrageSvcImpl) findChainsRecurse(ctx context.Context, currentAsset, targetAsset string, chain *domain.CandidateChain, chains *domain.CandidateChains, depth int) error {

	// create if nil
	if chain == nil {
		chain = &domain.CandidateChain{TotalRate: 1.0}
	}

	// apply restriction on maximum depth
	if depth >= s.cfg.Arbitrage.Depth {
		return nil
	}

	// request bids from provider
	bids, err := s.bidProvider.GetBidLightsBySourceAsset(ctx, currentAsset)
	if err != nil {
		return err
	}

	// go through bids and looking for possible conversions from the current asset
	var amount float64
	for _, r := range bids {

		// we aren't interested in rate 1
		if r.Rate == 1.0 {
			continue
		}

		if s.cfg.Arbitrage.CheckLimit {
			// skip chains which don't correspond minimum limits
			// we take prev amount here because limit is specified in the source asset
			if chain.Amount > 0.0 && chain.Amount < r.MinLimit {
				continue
			}

			// calc and check limits depending on available amounts in bid and amount from the previous bids
			if chain.Amount == 0.0 {
				// this is the first bid, so don't have previous amount
				amount = r.Available * r.Rate
			} else {
				// amount is calculates as min value of either prev amount converter to the current asset
				// or available amount of the current asset
				amount = math.Min(chain.Amount*r.Rate, r.Available)
			}
		}

		ch := s.copyChain(chain)
		ch.TotalRate = chain.TotalRate * r.Rate
		ch.Amount = amount
		ch.BidIds = append(ch.BidIds, r.Id)

		// if we've reached the target asset and total rate is greater than profitable rate min limit, then add a new chain to result
		if r.TrgAsset == targetAsset {
			// check minimum profit
			if ch.TotalRate < s.cfg.Arbitrage.MinProfit {
				continue
			}
			chains.Chains = append(chains.Chains, ch)
		} else {
			// analyze further stages recursively
			err = s.findChainsRecurse(ctx, r.TrgAsset, targetAsset, ch, chains, depth+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *arbitrageSvcImpl) profitableChainGenId(bidIds []string) string {
	hash, _ := hashstructure.Hash(bidIds, hashstructure.FormatV2, nil)
	return strconv.FormatUint(hash, 10)
}

// buildProfitableChains converts candidate chains to profitable chains
func (s *arbitrageSvcImpl) buildProfitableChains(ctx context.Context, candidates []*domain.CandidateChain) ([]*domain.ProfitableChain, error) {
	l := s.l().C(ctx).Mth("calc-profit").Trc()

	if len(candidates) == 0 {
		return nil, nil
	}

	// gather bids Ids for all chains
	var bidIds []string
	for _, chain := range candidates {
		bidIds = append(bidIds, chain.BidIds...)
	}

	// apply distinct to get unique slice of ids
	bidIds = kit.Strings(bidIds).Distinct()

	// get full bids by ids from storage
	bidDetails, err := s.bidProvider.GetBidsByIds(ctx, bidIds)
	if err != nil {
		return nil, err
	}

	// build map
	bidMap := make(map[string]*domain.Bid, len(bidDetails))
	for _, b := range bidDetails {
		bidMap[b.Id] = b
	}

	// for each candidate build a profitable chain
	var profitableChains []*domain.ProfitableChain
	now := kit.Now()
	for _, candidate := range candidates {
		bidsCount := len(candidate.BidIds)
		bids := make([]*domain.Bid, bidsCount)
		var methods kit.Strings
		var bidAssets kit.Strings
		var exchangeCodes kit.Strings
		for i, bidId := range candidate.BidIds {
			bid, ok := bidMap[bidId]
			// turns out haven't found a full bid (e.g. the bid gone away already), so skip such candidate
			if !ok {
				break
			}
			bids[i] = bid
			methods = append(methods, bid.Methods...)
			bidAssets = append(bidAssets, bid.TrgAsset)
			exchangeCodes = append(exchangeCodes, bid.ExchangeCode)
			// if the last bid, add a profitable chain
			if i == bidsCount-1 {
				// build chain id
				chainId := s.profitableChainGenId(candidate.BidIds)
				// check if profitable chain already exists
				exists, err := s.chainStorage.ProfitableChainExists(ctx, chainId)
				if err != nil {
					return nil, err
				}
				if exists {
					l.TrcF("%s exists", chainId)
					break
				}
				bidAssets = append([]string{bids[i].TrgAsset}, bidAssets...)
				chain := &domain.ProfitableChain{
					Id:            s.profitableChainGenId(candidate.BidIds),
					Asset:         bids[i].TrgAsset,
					ProfitShare:   candidate.TotalRate,
					Methods:       methods.Distinct(),
					BidAssets:     bidAssets,
					Bids:          bids,
					Depth:         bidsCount,
					ExchangeCodes: exchangeCodes.Distinct(),
					CreatedAt:     now,
				}
				profitableChains = append(profitableChains, chain)
				l.DbgF("chain(%s): asset:%s; ", chain.Id, chain.Asset)
			}
		}
	}

	// remove duplication
	chMap := make(map[string]*domain.ProfitableChain)
	for _, ch := range profitableChains {
		chMap[ch.Id] = ch
	}
	profitableChains = []*domain.ProfitableChain{}
	for _, v := range chMap {
		profitableChains = append(profitableChains, v)
	}

	return profitableChains, nil
}

func (s *arbitrageSvcImpl) assetsProviderWorker(ctx context.Context, tick time.Duration) {

	goroutine.New().
		WithLogger(s.l().C(ctx).Mth("assets-provider-worker")).
		WithRetry(goroutine.Unrestricted).
		WithRetryDelay(time.Second*10).
		Go(ctx, func() {
			l := s.l().C(ctx).Mth("assets-provider-worker").Trc()
			ticker := time.NewTicker(tick)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					assets, err := s.bidProvider.GetAssets(ctx)
					if err != nil {
						s.l().C(ctx).Mth("assets-provider-worker").E(err).Err("get-assets")
						continue
					}
					l.DbgF("%+v", assets)
					for _, asset := range assets {
						s.assetsToCalculateChan <- asset
					}
				case <-ctx.Done():
					l.Inf("stop")
					return
				}
			}
		})
}

func (s *arbitrageSvcImpl) findChainsWorker(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		i := i
		goroutine.New().
			WithLogger(s.l().C(ctx).Mth("find-chains-worker")).
			WithRetry(goroutine.Unrestricted).
			WithRetryDelay(time.Second*10).
			Go(ctx, func() {
				l := s.l().C(ctx).Mth("find-chains-worker").F(log.FF{"workerId": i}).Trc()
				for {
					select {
					case asset := <-s.assetsToCalculateChan:
						// find chains
						l.DbgF("analyzing %s", asset)
						chains := &domain.CandidateChains{}
						err := s.findChainsRecurse(ctx, asset, asset, nil, chains, 0)
						if err != nil {
							l.E(err).Err("find chains")
							continue
						}
						// send further to calc profit
						if len(chains.Chains) > 0 {
							l.DbgF("chain candidates: %s, chains: %d", asset, len(chains.Chains))
							s.processProfitableChainsChan <- chains.Chains
						}
					case <-ctx.Done():
						l.Inf("stop")
						return
					}
				}
			})
	}
}

func (s *arbitrageSvcImpl) profitableChainsProcessWorker(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		i := i
		goroutine.New().
			WithLogger(s.l().C(ctx).Mth("calc-profit-worker")).
			WithRetry(goroutine.Unrestricted).
			WithRetryDelay(time.Second*10).
			Go(ctx, func() {
				l := s.l().C(ctx).Mth("calc-profit-worker").F(log.FF{"workerId": i}).Trc()
				for {
					select {
					case candidates := <-s.processProfitableChainsChan:
						// calc profit
						profitableChains, err := s.buildProfitableChains(ctx, candidates)
						if err != nil {
							l.E(err).Err("calc profit chains")
							continue
						}
						if len(profitableChains) > 0 {
							s.saveProfitableChainsChan <- profitableChains
						}
					case <-ctx.Done():
						l.Inf("stop")
						return
					}
				}
			})
	}
}

func (s *arbitrageSvcImpl) saveProfitableChainsWorker(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		i := i
		goroutine.New().
			WithLogger(s.l().C(ctx).Mth("find-chains-worker")).
			WithRetry(goroutine.Unrestricted).
			WithRetryDelay(time.Second*10).
			Go(ctx, func() {
				l := s.l().C(ctx).Mth("find-chains-worker").F(log.FF{"workerId": i}).Trc()
				for {
					select {
					case chains := <-s.saveProfitableChainsChan:
						// save to store
						err := s.chainStorage.SaveProfitableChains(ctx, chains)
						if err != nil {
							l.E(err).Err("save chains")
							continue
						}
						// send to pipeline further
						s.profitableChainsNotifyChan <- chains
					case <-ctx.Done():
						l.Inf("stop")
						return
					}
				}
			})
	}
}

func (s *arbitrageSvcImpl) profitableChainsNotifyWorker(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		i := i
		goroutine.New().
			WithLogger(s.l().C(ctx).Mth("chains-notify-worker")).
			WithRetry(goroutine.Unrestricted).
			WithRetryDelay(time.Second*10).
			Go(ctx, func() {
				l := s.l().C(ctx).Mth("chains-notify-worker").F(log.FF{"workerId": i}).Trc()
				for {
					select {
					case chains := <-s.profitableChainsNotifyChan:
						l.TrcF("chains: %d", len(chains))
						if s.notifier != nil {
							if err := s.notifier.Notify(ctx, chains); err != nil {
								s.l().C(ctx).Mth("chains-notify-worker").E(err).Err()
							}
						}
					case <-ctx.Done():
						l.Inf("stop")
						return
					}
				}
			})
	}
}

func (s *arbitrageSvcImpl) RunCalculationBackground(ctx context.Context) error {
	l := s.l().C(ctx).Mth("run-calc").Trc()

	// check running
	if s.running.Load() {
		return errors.ErrChainsCalculationAlreadyRun(ctx)
	}

	// specify cancelled context
	ctx, s.cancelFunc = context.WithCancel(ctx)
	s.running.Store(true)

	// run provider
	if err := s.bidProvider.Run(ctx); err != nil {
		return err
	}

	// run workers
	s.assetsProviderWorker(ctx, time.Duration(s.cfg.Arbitrage.ProcessAssetsPeriodSec)*time.Second)
	s.findChainsWorker(ctx, 3)
	s.profitableChainsProcessWorker(ctx, 3)
	s.saveProfitableChainsWorker(ctx, 3)
	s.profitableChainsNotifyWorker(ctx, 3)

	l.Inf("ok")

	return nil
}

func (s *arbitrageSvcImpl) StopCalculation(ctx context.Context) error {
	l := s.l().C(ctx).Mth("stop-calc").Trc()
	// cancel if running
	if s.cancelFunc != nil && s.running.Load() {
		// stop
		s.cancelFunc()
		s.running.Store(false)
		s.cancelFunc = nil
		// stop bid provider
		_ = s.bidProvider.Stop(ctx)
		l.Inf("ok")
	}
	return nil
}

func (s *arbitrageSvcImpl) GetProfitableChains(ctx context.Context, rq *domain.GetProfitableChainsRequest) (*domain.GetProfitableChainsResponse, error) {
	s.l().C(ctx).Mth("get-profitable-chains").Trc()
	if rq.Size <= 0 {
		rq.Size = 100
	}
	return s.chainStorage.GetProfitableChains(ctx, rq)
}

func (s *arbitrageSvcImpl) GetProfitableChain(ctx context.Context, chainId string) (*domain.ProfitableChain, error) {
	s.l().C(ctx).Mth("get-profitable-chain-details").Trc()
	return s.chainStorage.GetProfitableChain(ctx, chainId)
}
