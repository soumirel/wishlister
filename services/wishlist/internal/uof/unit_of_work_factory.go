package uof

import "github.com/soumirel/wishlister/services/wishlist/internal/repository"

type uofFactory struct {
	connFactory repository.ConnFactory
}

func NewUnitOfWorkFactory(connFactory repository.ConnFactory) *uofFactory {
	return &uofFactory{
		connFactory: connFactory,
	}
}

func (uf *uofFactory) NewUnitOfWork(useTx bool) UnitOfWork {
	return newUnitOfWork(useTx, uf.connFactory)
}
