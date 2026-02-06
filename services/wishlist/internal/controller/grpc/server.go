package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/soumirel/wishlister/api/proto/gen/go/wishlist"
	"github.com/soumirel/wishlister/services/wishlist/internal/domain/entity"
	useridentuc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user_identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcServer struct {
	pb.UnimplementedWishlistServiceServer

	userIdentityUc *useridentuc.UserIdentityUsecase
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
	ctx context.Context, req *pb.GetUserIdByExternalIdentityRequest,
) (
	*pb.GetWishlistsResponse, error,
) {
	return nil, status.Error(codes.Unimplemented, "method GetWishlists not implemented")
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

func newGrpcServer(userIdentityUc *useridentuc.UserIdentityUsecase) pb.WishlistServiceServer {
	return &grpcServer{
		userIdentityUc: userIdentityUc,
	}
}

const (
	grpcPort = "8081"
)

func StartGrpcServer(userIdentityUc *useridentuc.UserIdentityUsecase) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", grpcPort))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	srv := newGrpcServer(userIdentityUc)
	pb.RegisterWishlistServiceServer(s, srv)
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
	return nil
}
