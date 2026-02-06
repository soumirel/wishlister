package repository

import (
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/repository"
)

type repoFactory struct {
	q Querier
}

func NewRepositoryFactory(q Querier) repository.RepositoryFactory {
	return &repoFactory{
		q: q,
	}
}

func (f *repoFactory) UserRepository() repository.UserRepository {
	return newUserRepository(f.q)
}

func (f *repoFactory) WishRepository() repository.WishRepository {
	return newWishRepository(f.q)
}

func (f *repoFactory) WishlistRepository() repository.WishlistRepository {
	return newWishlistRepository(f.q)
}

func (f *repoFactory) WishlistPermissionRepository() repository.WishlistPermissionRepository {
	return newWishlistPersmissionRepository(f.q)
}

func (f *repoFactory) UserIdentityRepository() repository.UserIdentityRepository {
	return newUserIdentityRepository(f.q)
}
