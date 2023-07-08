package storage

import (
	"context"
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
	"time"
)

func (c *chainStorageImpl) toProfitableChainAero(chain *domain.ProfitableChain) aero.BinMap {
	det, _ := json.Marshal(chain.Bids)
	return aero.BinMap{
		"asset":          chain.Asset,
		"profit_share":   chain.ProfitShare,
		"depth":          chain.Depth,
		"methods":        chain.Methods,
		"bid_assets":     chain.BidAssets,
		"exchange_codes": chain.ExchangeCodes,
		"created_at":     chain.CreatedAt.UnixNano(),
		"bids":           det,
	}
}

func (c *chainStorageImpl) toProfitableChainDomain(ctx context.Context, chain *aero.Record) (*domain.ProfitableChain, error) {
	if chain == nil {
		return nil, nil
	}
	createdAtInt, err := aerospike.AsInt(ctx, chain.Bins, "created_at")
	if err != nil {
		return nil, err
	}
	r := &domain.ProfitableChain{
		Id: chain.Key.Value().String(),
	}
	r.Asset, err = aerospike.AsString(ctx, chain.Bins, "asset")
	if err != nil {
		return nil, err
	}
	r.ProfitShare, err = aerospike.AsFloat(ctx, chain.Bins, "profit_share")
	if err != nil {
		return nil, err
	}
	r.Methods, err = aerospike.AsStrings(ctx, chain.Bins, "methods")
	if err != nil {
		return nil, err
	}
	r.BidAssets, err = aerospike.AsStrings(ctx, chain.Bins, "bid_assets")
	if err != nil {
		return nil, err
	}
	r.Depth, err = aerospike.AsInt(ctx, chain.Bins, "depth")
	if err != nil {
		return nil, err
	}
	r.ExchangeCodes, err = aerospike.AsStrings(ctx, chain.Bins, "exchange_codes")
	if err != nil {
		return nil, err
	}
	r.CreatedAt = time.Unix(0, int64(createdAtInt))
	bidsb, err := aerospike.AsBytes(ctx, chain.Bins, "bids")
	if err != nil {
		return nil, err
	}
	if bidsb != nil {
		_ = json.Unmarshal(bidsb, &r.Bids)
	}
	return r, nil
}
