package grpc

import (
	"context"

	pb "github.com/soumirel/wishlister/api/proto/gen/go/wishlist"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/model"
	"github.com/soumirel/wishlister/services/telegram-bot/internal/repository/grpc/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type wishlistGRPC struct {
	stub pb.WishlistServiceClient
}

func NewWishlistGrpcRepository(addr string) (*wishlistGRPC, error) {
	var opts []grpc.DialOption

	opts = append(opts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10MB
			grpc.MaxCallSendMsgSize(10*1024*1024),
		),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithChainUnaryInterceptor(
			interceptors.AuthUnaryInterceptor,
		),
	)

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	client := &wishlistGRPC{
		stub: pb.NewWishlistServiceClient(conn),
	}

	return client, nil
}

func (r *wishlistGRPC) GetUserIdByExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error) {
	req := pb.GetUserIdByExternalIdentityRequest{
		ExternalIdentity: &pb.ExternalIdentity{
			ExternalID:       ei.ExternalID,
			IdentityProvider: ei.IdentityProvider,
		},
	}
	resp, err := r.stub.GetUserIdByExternalIdentity(ctx, &req)
	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			return "", model.ErrUserDoesNotExist

		}
		return "", err
	}
	userID := resp.GetUserID()
	return userID, nil
}

func (r *wishlistGRPC) CreateUserFromExternalIdentity(ctx context.Context, ei model.ExternalIdentity) (string, error) {
	req := pb.CreateUserFromExternalIdentityRequest{
		ExternalIdentity: &pb.ExternalIdentity{
			ExternalID:       ei.ExternalID,
			IdentityProvider: ei.IdentityProvider,
		},
	}
	resp, err := r.stub.CreateUserFromExternalIdentity(ctx, &req)
	if err != nil {
		return "", err
	}
	userID := resp.GetUserID()
	return userID, nil
}

func (r *wishlistGRPC) GetWishlists(ctx context.Context) (model.WishlistList, error) {
	req := pb.GetWishlistsRequest{}
	resp, err := r.stub.GetWishlists(ctx, &req)
	if err != nil {
		return nil, err
	}
	wishlistsResp := resp.GetWishlists()
	list := make(model.WishlistList, len(wishlistsResp))
	for i, v := range wishlistsResp {
		list[i] = &model.Wishlist{
			ID:     v.ID,
			UserID: v.UserID,
			Name:   v.Name,
		}
	}
	return list, nil
}

func (r *wishlistGRPC) CreateWishlist(ctx context.Context, w model.Wishlist) (model.Wishlist, error) {
	req := pb.CreateWishlistRequest{
		Name: w.Name,
	}
	resp, err := r.stub.CreateWishlist(ctx, &pb.CreateWishlistRequest{
		Name: req.Name,
	})
	if err != nil {
		return model.Wishlist{}, err
	}
	wishlistResp := resp.GetWishlist()
	wishlist := model.Wishlist{
		ID:     wishlistResp.GetID(),
		Name:   wishlistResp.GetName(),
		UserID: wishlistResp.GetUserID(),
	}
	return wishlist, nil
}
