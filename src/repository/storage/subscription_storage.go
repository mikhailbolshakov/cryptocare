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
	SetSubscriptions = "subscriptions"
)

type subscriptionDetails struct {
	Notifications []*domain.SubscriptionNotification `json:"notifications,omitempty"` // Notifications notifications
}

type subscriptionStorageImpl struct {
	aero kitAero.Aerospike
	cfg  *kitAero.Config
}

func (s *subscriptionStorageImpl) l() log.CLogger {
	return service.L().Cmp("subscription-storage")
}

func newSubscriptionStorage(aero kitAero.Aerospike, cfg *kitAero.Config) *subscriptionStorageImpl {
	return &subscriptionStorageImpl{
		aero: aero,
		cfg:  cfg,
	}
}

func (s *subscriptionStorageImpl) SaveSubscription(ctx context.Context, subs *domain.Subscription) error {
	s.l().C(ctx).Mth("save").Trc()
	writePolicy := aero.NewWritePolicy(0, 0)
	writePolicy.SendKey = true
	key, err := aero.NewKey(s.cfg.Namespace, SetSubscriptions, subs.Id)
	if err != nil {
		return errors.ErrSubscriptionStoragePut(err, ctx)
	}
	err = s.aero.Instance().Put(writePolicy, key, s.toSubscriptionAero(subs))
	if err != nil {
		return errors.ErrSubscriptionStoragePut(err, ctx)
	}
	return nil
}

func (s *subscriptionStorageImpl) GetSubscription(ctx context.Context, subsId string) (*domain.Subscription, error) {
	s.l().C(ctx).Mth("get").F(log.FF{"subscriptionId": subsId}).Trc()

	key, aeroErr := aero.NewKey(s.cfg.Namespace, SetSubscriptions, subsId)
	if aeroErr != nil {
		return nil, errors.ErrSubscriptionStorageGet(aeroErr, ctx)
	}

	policy := aero.NewPolicy()
	policy.SendKey = true
	rec, aeroErr := s.aero.Instance().Get(policy, key)
	if aeroErr != nil && !aeroErr.Matches(types.KEY_NOT_FOUND_ERROR) {
		return nil, errors.ErrSubscriptionStorageGet(aeroErr, ctx)
	}

	subs, err := s.toSubscriptionDomain(ctx, rec)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *subscriptionStorageImpl) DeleteSubscription(ctx context.Context, subsId string) error {
	s.l().Mth("delete").C(ctx).F(log.FF{"subsId": subsId}).Trc()
	key, err := aero.NewKey(s.cfg.Namespace, SetSubscriptions, subsId)
	if err != nil {
		return errors.ErrSubscriptionStorageDel(err, ctx)
	}
	_, err = s.aero.Instance().Delete(nil, key)
	if err != nil {
		return errors.ErrSubscriptionStorageDel(err, ctx)
	}
	return nil
}

func (s *subscriptionStorageImpl) SearchSubscriptions(ctx context.Context, rq *domain.SearchSubscriptionsRequest) ([]*domain.Subscription, error) {
	s.l().C(ctx).Mth("search").Trc()

	exp := aero.ExpEq(aero.ExpBoolVal(true), aero.ExpBoolVal(true))
	if !rq.WithInActive {
		exp = aero.ExpAnd(exp, aero.ExpEq(aero.ExpBoolBin("is_active"), aero.ExpBoolVal(true)))
	}
	if rq.UserId != "" {
		exp = aero.ExpAnd(exp, aero.ExpEq(aero.ExpStringBin("user_id"), aero.ExpStringVal(rq.UserId)))
	}

	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.SendKey = true
	queryPolicy.FilterExpression = exp
	statement := aero.NewStatement(s.cfg.Namespace, SetSubscriptions)

	recordSet, err := s.aero.Instance().Query(queryPolicy, statement)
	if err != nil {
		return nil, errors.ErrSubscriptionStorageSearch(err, ctx)
	}
	var res []*domain.Subscription
	for r := range recordSet.Results() {
		if r.Err != nil {
			return nil, errors.ErrSubscriptionStorageSearch(r.Err, ctx)
		} else {
			subs, err := s.toSubscriptionDomain(ctx, r.Record)
			if err != nil {
				return nil, err
			}
			res = append(res, subs)
		}
	}
	return res, nil
}
