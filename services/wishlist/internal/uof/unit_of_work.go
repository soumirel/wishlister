package uof

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	domainRepo "github.com/soumirel/wishlister/wishlist/internal/domain/repository"
	"github.com/soumirel/wishlister/wishlist/internal/repository"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, rf domainRepo.RepositoryFactory) error) error
}

type unitOfWork struct {
	done     *atomic.Bool
	initOnce *sync.Once

	useTx bool

	conn repository.Conn
	tx   repository.Tx

	connFactory repository.ConnFactory
}

func newUnitOfWork(
	useTx bool, connFactory repository.ConnFactory,
) *unitOfWork {
	return &unitOfWork{
		useTx:       useTx,
		done:        &atomic.Bool{},
		initOnce:    &sync.Once{},
		connFactory: connFactory,
	}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, rf domainRepo.RepositoryFactory) error) error {
	if u.done.Load() {
		return errors.New("unit of work done")
	}
	defer u.done.Store(true)

	var err error
	u.initOnce.Do(func() {
		u.conn, err = u.connFactory.GetConn(ctx)
		if err != nil {
			return
		}
		if !u.useTx {
			return
		}
		u.tx, err = u.conn.Begin(ctx)
		if err != nil {
			u.conn.Release()
			return
		}
	})
	if err != nil {
		return err
	}
	if (u.conn == nil) || (u.useTx && u.tx == nil) {
		return errors.New("connection init failed")
	}

	defer func() {
		if u.conn != nil {
			u.conn.Release()
			u.conn = nil
			u.tx = nil
		}
		u.done.Store(true)
	}()

	var querier repository.Querier
	if u.useTx {
		querier = u.tx
	} else {
		querier = u.conn
	}
	repoFactory := repository.NewRepositoryFactory(querier)
	err = fn(ctx, repoFactory)
	if u.tx != nil {
		if err != nil {
			err := u.tx.Rollback(ctx)
			if err != nil {
				log.Println(fmt.Errorf("rollback error: %w", err))
			}
		} else {
			err := u.tx.Commit(ctx)
			if err != nil {
				return err
			}
		}
	}

	return err
}
