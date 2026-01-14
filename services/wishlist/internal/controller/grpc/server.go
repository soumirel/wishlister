package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/soumirel/wishlister/api/proto/wishlist"
	useridentuc "github.com/soumirel/wishlister/wishlist/internal/usecase/user_identity"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedWishlistServer

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
		return nil, err
	}
	resp := pb.GetUserIdByExternalIdentityResponse{
		UserID: userID,
	}
	return &resp, nil
}

func newGrpcServer(userIdentityUc *useridentuc.UserIdentityUsecase) pb.WishlistServer {
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
	pb.RegisterWishlistServer(s, srv)
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
	return nil
}
