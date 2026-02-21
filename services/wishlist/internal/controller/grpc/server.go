package grpc

import (
	"context"
	"errors"
	"net"

	pb "github.com/soumirel/wishlister/api/proto/gen/go/wishlist"
	"github.com/soumirel/wishlister/services/wishlist/internal/auth"
	"github.com/soumirel/wishlister/services/wishlist/internal/controller/grpc/interceptors"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	useridentuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user_identity"
	wishlistuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/wishlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	pb.UnimplementedWishlistServiceServer

	userIdentityUc *useridentuc.UserIdentityUsecase
	wishlistUc     *wishlistuc.WishlistUsecase
}

func newGrpcServer(
	userIdentityUc *useridentuc.UserIdentityUsecase,
	wishlistUc *wishlistuc.WishlistUsecase,
) pb.WishlistServiceServer {
	return &grpcServer{
		userIdentityUc: userIdentityUc,
		wishlistUc:     wishlistUc,
	}
}

func StartGrpcServer(
	grpcAddr string,
	userIdentityUc *useridentuc.UserIdentityUsecase,
	wishlistUc *wishlistuc.WishlistUsecase,
) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.AuthUnaryInterceptor),
	}
	s := grpc.NewServer(opts...)
	srv := newGrpcServer(
		userIdentityUc, wishlistUc,
	)
	pb.RegisterWishlistServiceServer(s, srv)
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *grpcServer) CreateUserFromExternalIdentity(
	ctx context.Context, req *pb.CreateUserFromExternalIdentityRequest) (
	*pb.CreateUserFromExternalIdentityResponse, error,
) {
	cmd := useridentuc.CreateUserFromExternalIdentityCommand{
		ExternalID:       req.GetExternalIdentity().GetExternalID(),
		IdentityProvider: req.GetExternalIdentity().GetIdentityProvider(),
	}
	userID, err := s.userIdentityUc.CreateUserFromExternalIdentity(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateUserFromExternalIdentityResponse{
		UserID: userID,
	}
	return &resp, nil
}

func (s *grpcServer) GetWishlists(
	ctx context.Context, req *pb.GetWishlistsRequest,
) (
	*pb.GetWishlistsResponse, error,
) {
	au, ok := auth.FromCtx(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}
	cmd := wishlistuc.GetWishlistsCommand{
		RequestorUserID: au.UserID,
	}
	wishlists, err := s.wishlistUc.GetWishlists(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := pb.GetWishlistsResponse{
		Wishlists: make([]*pb.Wishlist, len(wishlists)),
	}
	for i, w := range wishlists {
		resp.Wishlists[i] = &pb.Wishlist{
			ID:     w.ID,
			UserID: w.UserID,
			Name:   w.Name,
		}
	}
	return &resp, nil
}

func (s *grpcServer) GetUserIdByExternalIdentity(
	ctx context.Context, req *pb.GetUserIdByExternalIdentityRequest,
) (
	*pb.GetUserIdByExternalIdentityResponse, error,
) {
	cmd := useridentuc.GetUserIdByExternalIdentityCommand{
		ExternalID:       req.GetExternalIdentity().GetExternalID(),
		IdentityProvider: req.GetExternalIdentity().GetIdentityProvider(),
	}
	userID, err := s.userIdentityUc.GetUserIdByExternalIdentity(ctx, cmd)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserIdentityDoesNotExist):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	resp := pb.GetUserIdByExternalIdentityResponse{
		UserID: userID,
	}
	return &resp, nil
}

func (s *grpcServer) CreateWishlist(
	ctx context.Context, req *pb.CreateWishlistRequest,
) (*pb.CreateWishlistResponse, error) {
	a, ok := auth.FromCtx(ctx)
	if !ok {
		return nil, auth.ErrUnauthorized
	}
	cmd := wishlistuc.CreateWishlistCommand{
		RequestorUserID: a.UserID,
		Name:            req.GetName(),
	}
	wishlist, err := s.wishlistUc.CreateWishlist(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateWishlistResponse{
		Wishlist: &pb.Wishlist{
			ID:     wishlist.ID,
			UserID: wishlist.UserID,
			Name:   wishlist.Name,
		},
	}
	return &resp, nil
}
