package server

import (
	"context"
	"github.com/shewitt93/explore_service/internal/entity"
	"github.com/shewitt93/explore_service/internal/repository"
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExploreGRPCServer struct {
	grpclibs.ExploreServiceServer
	repo repository.DecisionRepository
}

func NewExploreGRPCServer(repo repository.DecisionRepository) *ExploreGRPCServer {
	return &ExploreGRPCServer{
		repo: repo,
	}
}

func (s *ExploreGRPCServer) ListLikedYou(ctx context.Context, req *grpclibs.ListLikedYouRequest) (*grpclibs.ListLikedYouResponse, error) {

	//user isn't empty check if receipient ID is not empty string
	if req.GetRecipientUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing recipient user id")
	}
	var cursor *entity.Cursor
	if req.PaginationToken != nil && *req.PaginationToken != "" {
		decodedCursor, err := entity.DecodeCursor(*req.PaginationToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pagination token: %v", err)
		}
		cursor = decodedCursor
	}

	const pageSize = 50

	likers, nextCursor, err := s.repo.ListLikersByRecipient(ctx, req.GetRecipientUserId(), cursor, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch likers: %v", err)
	}

	response := &grpclibs.ListLikedYouResponse{
		Likers: make([]*grpclibs.ListLikedYouResponse_Liker, 0, len(likers)),
	}

	for _, liker := range likers {
		response.Likers = append(response.Likers, &grpclibs.ListLikedYouResponse_Liker{
			ActorId:       liker.ActorID,
			UnixTimestamp: liker.UnixTimestamp,
		})
	}

	if nextCursor != nil {
		token, err := entity.EncodeCursor(nextCursor)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to encode pagination token: %v", err)
		}
		response.NextPaginationToken = &token
	}

	return response, nil
}

func (s *ExploreGRPCServer) ListNewLikedYou(ctx context.Context, req *grpclibs.ListLikedYouRequest) (*grpclibs.ListLikedYouResponse, error) {

	// Check if recipient ID is not empty string
	if req.GetRecipientUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing recipient user id")
	}

	// Handle pagination token if provided
	var cursor *entity.Cursor
	if req.PaginationToken != nil && *req.PaginationToken != "" {
		decodedCursor, err := entity.DecodeCursor(*req.PaginationToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pagination token: %v", err)
		}
		cursor = decodedCursor
	}

	// Define page size
	const pageSize = 50

	// Call repository function to fetch new likers
	likers, nextCursor, err := s.repo.ListNewLikersByRecipient(ctx, req.GetRecipientUserId(), cursor, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch new likers: %v", err)
	}

	// Build response
	response := &grpclibs.ListLikedYouResponse{
		Likers: make([]*grpclibs.ListLikedYouResponse_Liker, 0, len(likers)),
	}

	// Add likers to response
	for _, liker := range likers {
		response.Likers = append(response.Likers, &grpclibs.ListLikedYouResponse_Liker{
			ActorId:       liker.ActorID,
			UnixTimestamp: liker.UnixTimestamp,
		})
	}

	// Add pagination token if there are more results
	if nextCursor != nil {
		token, err := entity.EncodeCursor(nextCursor)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to encode pagination token: %v", err)
		}
		response.NextPaginationToken = &token
	}

	return response, nil
}

func (s *ExploreGRPCServer) CountLikedYou(ctx context.Context, req *grpclibs.CountLikedYouRequest) (*grpclibs.CountLikedYouResponse, error) {

	if req.GetRecipientUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing recipient user id")
	}

	count, err := s.repo.CountLikersByRecipient(ctx, req.GetRecipientUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count likers: %v", err)
	}

	return &grpclibs.CountLikedYouResponse{
		Count: count,
	}, nil
}

func (s *ExploreGRPCServer) PutDecision(ctx context.Context, req *grpclibs.PutDecisionRequest) (*grpclibs.PutDecisionResponse, error) {
	// Validate input
	if req.GetActorUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing actor user id")
	}
	if req.GetRecipientUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing recipient user id")
	}

	// Call repository function to put decision
	mutualLike, err := s.repo.CreateOrUpdateDecision(ctx, req.GetActorUserId(), req.GetRecipientUserId(), req.GetLikedRecipient())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to put decision: %v", err)
	}

	// Return response with mutual like status
	return &grpclibs.PutDecisionResponse{
		MutualLikes: mutualLike,
	}, nil
}
