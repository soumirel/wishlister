package usecase

import "wishlister/internal/uof"

type UnitOfWorkFactory interface {
	NewUnitOfWork(useTx bool) uof.UnitOfWork
}
