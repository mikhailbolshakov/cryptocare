package storage

import (
	"context"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
)

func (b *bidStorageImpl) toBidLightDomain(ctx context.Context, dto *aero.Record) (*domain.BidLight, error) {
	if dto == nil {
		return nil, nil
	}
	r := &domain.BidLight{Id: dto.Key.Value().String(), Type: domain.BidTypeP2P}
	var err error
	r.SrcAsset, err = aerospike.AsString(ctx, dto.Bins, "src")
	if err != nil {
		return nil, err
	}
	r.TrgAsset, err = aerospike.AsString(ctx, dto.Bins, "trg")
	if err != nil {
		return nil, err
	}
	r.Rate, err = aerospike.AsFloat(ctx, dto.Bins, "rate")
	if err != nil {
		return nil, err
	}
	r.Available, err = aerospike.AsFloat(ctx, dto.Bins, "available")
	if err != nil {
		return nil, err
	}
	r.MinLimit, err = aerospike.AsFloat(ctx, dto.Bins, "minLimit")
	if err != nil {
		return nil, err
	}
	r.MaxLimit, err = aerospike.AsFloat(ctx, dto.Bins, "maxLimit")
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (b *bidStorageImpl) toBidDomain(ctx context.Context, dto *aero.Record) (*domain.Bid, error) {
	if dto == nil {
		return nil, nil
	}
	r := &domain.Bid{Id: dto.Key.Value().String(), Type: domain.BidTypeP2P}
	var err error
	r.SrcAsset, err = aerospike.AsString(ctx, dto.Bins, "src")
	if err != nil {
		return nil, err
	}
	r.TrgAsset, err = aerospike.AsString(ctx, dto.Bins, "trg")
	if err != nil {
		return nil, err
	}
	r.Rate, err = aerospike.AsFloat(ctx, dto.Bins, "rate")
	if err != nil {
		return nil, err
	}
	r.Available, err = aerospike.AsFloat(ctx, dto.Bins, "available")
	if err != nil {
		return nil, err
	}
	r.MinLimit, err = aerospike.AsFloat(ctx, dto.Bins, "minLimit")
	if err != nil {
		return nil, err
	}
	r.MaxLimit, err = aerospike.AsFloat(ctx, dto.Bins, "maxLimit")
	if err != nil {
		return nil, err
	}
	r.UserId, err = aerospike.AsString(ctx, dto.Bins, "userId")
	if err != nil {
		return nil, err
	}
	r.ExchangeCode, err = aerospike.AsString(ctx, dto.Bins, "exchangeCode")
	if err != nil {
		return nil, err
	}
	r.Link, err = aerospike.AsString(ctx, dto.Bins, "link")
	if err != nil {
		return nil, err
	}
	r.Methods, err = aerospike.AsStrings(ctx, dto.Bins, "methods")
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (b *bidStorageImpl) toBidAero(bid *domain.Bid) aero.BinMap {
	return aero.BinMap{
		"src":          bid.SrcAsset,
		"trg":          bid.TrgAsset,
		"rate":         bid.Rate,
		"exchangeCode": bid.ExchangeCode,
		"available":    bid.Available,
		"minLimit":     bid.MinLimit,
		"maxLimit":     bid.MaxLimit,
		"methods":      bid.Methods,
		"userId":       bid.UserId,
		"link":         bid.Link,
	}
}
