package server

import (
	"context"
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ExploreGRPCServer struct {
	grpclibs.ExploreServiceServer
}

func NewExploreGRPCServer() *ExploreGRPCServer {
	return &ExploreGRPCServer{}
}

func (s *ExploreGRPCServer) ListLikedYou(ctx context.Context, req *grpclibs.ListLikedYouRequest) (*grpclibs.ListLikedYouResponse, error) {

	//user isn't empty check if receipient ID is not empty string
	if req.GetRecipientUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing recipient user id")
	}

	//paginate on recip id and updated at as cursor

	//check if pagination token is present otherwise return all

	return
}

func (s *ExploreGRPCServer) ListNewLikedYou(ctx context.Context req.grpclibs.ListNewLikedYouRequest) (*grpclibs.ListLikedYouResponse, error) {
	return
}
