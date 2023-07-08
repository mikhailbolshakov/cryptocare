package storage

import (
	"context"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/aerospike/aerospike-client-go/v6/types"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	kitAero "github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
	"github.com/mikhailbolshakov/cryptocare/src/service"
)

const (
	SetProfitableChains = "profitable_chains"
)

type chainStorageImpl struct {
	aero kitAero.Aerospike
	cfg  *kitAero.Config
}

func (c *chainStorageImpl) l() log.CLogger {
	return service.L().Cmp("chain-storage")
}

func newChainStorage(aero kitAero.Aerospike, cfg *kitAero.Config) *chainStorageImpl {
	return &chainStorageImpl{
		aero: aero,
		cfg:  cfg,
	}
}

func (c *chainStorageImpl) SaveProfitableChains(ctx context.Context, chains []*domain.ProfitableChain) error {
	c.l().C(ctx).Mth("save-chains").Trc()
	writePolicy := aero.NewWritePolicy(0, 60*60)
	writePolicy.SendKey = true
	for _, chain := range chains {
		key, err := aero.NewKey(c.cfg.Namespace, SetProfitableChains, chain.Id)
		if err != nil {
			return errors.ErrChainStoragePutChain(err, ctx)
		}
		err = c.aero.Instance().Put(writePolicy, key, c.toProfitableChainAero(chain))
		if err != nil {
			return errors.ErrChainStoragePutChain(err, ctx)
		}
	}
	return nil
}

func (c *chainStorageImpl) GetProfitableChains(ctx context.Context, rq *domain.GetProfitableChainsRequest) (*domain.GetProfitableChainsResponse, error) {
	c.l().C(ctx).Mth("get-chains").Trc()

	var exp *aero.Expression

	// filter by assets
	var assetExps []*aero.Expression
	for _, asset := range rq.Assets {
		assetExps = append(assetExps, aero.ExpEq(aero.ExpStringBin("asset"), aero.ExpStringVal(asset)))
	}
	if len(assetExps) > 1 {
		exp = aero.ExpOr(assetExps...)
	}
	if len(assetExps) == 1 {
		exp = assetExps[0]
	}

	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.SendKey = true
	queryPolicy.MaxRecords = int64(rq.Size)
	queryPolicy.FilterExpression = exp

	bins := []string{"asset", "profit_share", "methods", "bid_assets", "depth", "exchange_codes", "created_at"}
	if rq.WithBids {
		bins = append(bins, "bids")
	}
	statement := aero.NewStatement(c.cfg.Namespace, SetProfitableChains, bins...)

	recordSet, err := c.aero.Instance().Query(queryPolicy, statement)
	if err != nil {
		return nil, errors.ErrChainStorageScanChains(err, ctx)
	}
	res := &domain.GetProfitableChainsResponse{}
	for r := range recordSet.Results() {
		if r.Err != nil {
			return nil, errors.ErrChainStorageScanChains(r.Err, ctx)
		} else {
			chain, err := c.toProfitableChainDomain(ctx, r.Record)
			if err != nil {
				return nil, err
			}
			res.Chains = append(res.Chains, chain)
		}
	}
	return res, nil
}

func (c *chainStorageImpl) GetProfitableChain(ctx context.Context, chainId string) (*domain.ProfitableChain, error) {
	c.l().C(ctx).Mth("get-chain").F(log.FF{"chainId": chainId}).Trc()

	key, aeroErr := aero.NewKey(c.cfg.Namespace, SetProfitableChains, chainId)
	if aeroErr != nil {
		return nil, errors.ErrChainStorageGetChain(aeroErr, ctx)
	}

	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, aeroErr := c.aero.Instance().Get(policy, key)
	if aeroErr != nil && !aeroErr.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrChainStorageGetChain(aeroErr, ctx)
	}

	chain, err := c.toProfitableChainDomain(ctx, rec)
	if err != nil {
		return nil, err
	}
	return chain, nil
}

func (c *chainStorageImpl) ProfitableChainExists(ctx context.Context, chainId string) (bool, error) {
	c.l().C(ctx).Mth("chain-exists").F(log.FF{"chainId": chainId}).Trc()

	key, err := aero.NewKey(c.cfg.Namespace, SetProfitableChains, chainId)
	if err != nil {
		return false, errors.ErrChainStorageGetChain(err, ctx)
	}

	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, err := c.aero.Instance().Get(policy, key)
	if err != nil {
		if err.Matches(types.KEY_NOT_FOUND_ERROR) {
			return false, nil
		}
		return false, errors.ErrChainStorageGetChain(err, ctx)
	}
	return rec != nil, nil
}
