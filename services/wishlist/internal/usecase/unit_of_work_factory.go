package usecase

import "github.com/soumirel/wishlister/wishlist/internal/uof"

type UnitOfWorkFactory interface {
	NewUnitOfWork(useTx bool) uof.UnitOfWork
}
