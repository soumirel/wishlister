package usecase

import "github.com/soumirel/wishlister/services/wishlist/internal/uof"

type UnitOfWorkFactory interface {
	NewUnitOfWork(useTx bool) uof.UnitOfWork
}
