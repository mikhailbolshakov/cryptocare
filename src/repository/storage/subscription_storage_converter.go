package storage

import (
	"context"
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
)

func (s *subscriptionStorageImpl) toSubscriptionAero(subs *domain.Subscription) aero.BinMap {
	det, _ := json.Marshal(&subscriptionDetails{
		Notifications: subs.Notifications,
	})
	return aero.BinMap{
		"user_id":        subs.UserId,
		"is_active":      subs.IsActive,
		"flt_assets":     subs.Filter.Assets,
		"flt_methods":    subs.Filter.Methods,
		"flt_exchanges":  subs.Filter.Exchanges,
		"flt_min_profit": subs.Filter.MinProfit,
		"flt_max_depth":  subs.Filter.MaxDepth,
		"details":        det,
	}
}

func (s *subscriptionStorageImpl) toSubscriptionDomain(ctx context.Context, subs *aero.Record) (*domain.Subscription, error) {
	if subs == nil {
		return nil, nil
	}
	r := &domain.Subscription{
		Id:     subs.Key.Value().String(),
		Filter: &domain.SubscriptionChainFilter{},
	}
	var err error
	r.UserId, err = aerospike.AsString(ctx, subs.Bins, "user_id")
	if err != nil {
		return nil, err
	}
	r.IsActive, err = aerospike.AsBool(ctx, subs.Bins, "is_active")
	if err != nil {
		return nil, err
	}
	r.Filter.Exchanges, err = aerospike.AsStrings(ctx, subs.Bins, "flt_exchanges")
	if err != nil {
		return nil, err
	}
	r.Filter.Assets, err = aerospike.AsStrings(ctx, subs.Bins, "flt_assets")
	if err != nil {
		return nil, err
	}
	r.Filter.Methods, err = aerospike.AsStrings(ctx, subs.Bins, "flt_methods")
	if err != nil {
		return nil, err
	}
	r.Filter.MaxDepth, err = aerospike.AsInt(ctx, subs.Bins, "flt_max_depth")
	if err != nil {
		return nil, err
	}
	r.Filter.MinProfit, err = aerospike.AsFloat(ctx, subs.Bins, "flt_min_profit")
	if err != nil {
		return nil, err
	}
	details, err := aerospike.AsBytes(ctx, subs.Bins, "details")
	if err != nil {
		return nil, err
	}
	if details != nil {
		det := &subscriptionDetails{}
		_ = json.Unmarshal(details, &det)
		r.Notifications = det.Notifications
	}
	return r, nil
}
