package storage

import (
	"context"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	kitAero "github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
	"github.com/mikhailbolshakov/cryptocare/src/service"
)

const (
	SetBidsP2P = "bids_p2p"
)

type bidStorageImpl struct {
	aero kitAero.Aerospike
	cfg  *kitAero.Config
}

func (b *bidStorageImpl) l() log.CLogger {
	return service.L().Cmp("bid-storage")
}

func newBidStorage(aero kitAero.Aerospike, cfg *kitAero.Config) *bidStorageImpl {
	return &bidStorageImpl{
		aero: aero,
		cfg:  cfg,
	}
}

func (b *bidStorageImpl) GetBidsByIds(ctx context.Context, ids []string) ([]*domain.Bid, error) {
	b.l().C(ctx).Mth("get-bids-by-ids").Trc()
	// build keys
	keys := make([]*aero.Key, len(ids))
	for i, id := range ids {
		key, _ := aero.NewKey(b.cfg.Namespace, SetBidsP2P, id)
		keys[i] = key
	}
	batchPolicy := aero.NewBatchPolicy()
	records, err := b.aero.Instance().BatchGet(batchPolicy, keys)
	if err != nil {
		return nil, errors.ErrBidStorageGetBidsByIds(err, ctx)
	}
	//res := make([]*domain.Bid, len(records))
	var res []*domain.Bid
	for _, r := range records {
		if r != nil {
			bd, err := b.toBidDomain(ctx, r)
			if err != nil {
				return nil, err
			}
			res = append(res, bd)
		}
	}
	return res, nil
}

func (b *bidStorageImpl) PutBids(ctx context.Context, bids []*domain.Bid, ttlSec uint32) error {
	b.l().C(ctx).Mth("put-bids").Trc()
	writePolicy := aero.NewWritePolicy(0, ttlSec)
	writePolicy.SendKey = true
	for _, bid := range bids {
		key, err := aero.NewKey(b.cfg.Namespace, SetBidsP2P, bid.Id)
		if err != nil {
			return errors.ErrBidStoragePutBids(err, ctx)
		}
		err = b.aero.Instance().Put(writePolicy, key, b.toBidAero(bid))
		if err != nil {
			return errors.ErrBidStoragePutBids(err, ctx)
		}
	}
	return nil
}

func (b *bidStorageImpl) GetBidsLightAll(ctx context.Context) ([]*domain.BidLight, error) {
	l := b.l().C(ctx).Mth("get-bids-light-all").Trc()
	// scan all bids
	scanPolicy := aero.NewScanPolicy()
	recordSet, err := b.aero.Instance().ScanAll(scanPolicy, b.cfg.Namespace, SetBidsP2P,
		"src", "trg", "rate", "minLimit", "maxLimit", "available")
	if err != nil {
		return nil, errors.ErrBidStorageScanBidsLight(err, ctx)
	}
	var res []*domain.BidLight
	for r := range recordSet.Results() {
		if r.Err != nil {
			l.E(errors.ErrBidStorageScanBidsLightReadRec(r.Err, ctx)).Err()
		} else {
			bd, err := b.toBidLightDomain(ctx, r.Record)
			if err != nil {
				return nil, err
			}
			res = append(res, bd)
		}
	}
	return res, nil
}
