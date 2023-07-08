package storage

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	kitService "github.com/mikhailbolshakov/cryptocare/src/kit/service"
	kitAero "github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/pg"
	"github.com/mikhailbolshakov/cryptocare/src/service"
)

type Adapter interface {
	kitService.StorageAdapter
	domain.BidStorage
	domain.ChainStorage
	domain.UserStorage
	domain.SubscriptionStorage
	auth.SessionStorage
}

type adapterImpl struct {
	*bidStorageImpl
	*chainStorageImpl
	*userStorageImpl
	*sessionStorageImpl
	*subscriptionStorageImpl
	aero kitAero.Aerospike
	pg   *pg.Storage
}

func NewAdapter() Adapter {
	a := &adapterImpl{
		aero: kitAero.New(),
	}
	return a
}

func (c *adapterImpl) Init(ctx context.Context, cfg interface{}) error {
	config := cfg.(*service.Config)

	// init postgres
	var err error
	c.pg, err = pg.Open(config.Storages.Pg.Master, service.LF())
	if err != nil {
		return err
	}

	// applying migrations
	if config.Storages.Pg.MigPath != "" {
		db, _ := c.pg.Instance.DB()
		m := pg.NewMigration(db, config.Storages.Pg.MigPath, service.LF())
		if err := m.Up(); err != nil {
			return err
		}
	}

	// init aero
	err = c.aero.Open(ctx, config.Storages.Aero, service.LF())
	if err != nil {
		return err
	}

	// init storages
	c.bidStorageImpl = newBidStorage(c.aero, config.Storages.Aero)
	c.chainStorageImpl = newChainStorage(c.aero, config.Storages.Aero)
	c.userStorageImpl = newUserStorage(c.pg, c.aero, config.Storages.Aero)
	c.subscriptionStorageImpl = newSubscriptionStorage(c.aero, config.Storages.Aero)
	err = c.userStorageImpl.init(ctx)
	if err != nil {
		return err
	}
	c.sessionStorageImpl = newSessionStorage(c.pg, c.aero, config.Storages.Aero)
	return nil
}

func (c *adapterImpl) Close(ctx context.Context) error {
	if c.aero != nil {
		_ = c.aero.Close(ctx)
	}
	if c.pg != nil {
		c.pg.Close()
	}
	return nil
}
