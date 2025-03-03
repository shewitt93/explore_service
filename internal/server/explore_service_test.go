package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shewitt93/explore_service/internal/entity"
	"github.com/shewitt93/explore_service/internal/repository"
	"github.com/shewitt93/explore_service/pkg/grpclibs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListLikedYou(t *testing.T) {
	// Create a context for testing
	ctx := context.Background()

	t.Run("Success_WithoutCursor", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		pageSize := 50
		likers := []entity.Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
			{ActorID: "actor2", UnixTimestamp: 1738686000},
		}

		// Set up mock expectations
		mockRepo.On("ListLikersByRecipient", ctx, recipientID, (*entity.Cursor)(nil), pageSize).
			Return(likers, nil, nil)

		// Create server with mock repository and cursor service
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Likers, 2)
		assert.Equal(t, "actor1", resp.Likers[0].ActorId)
		assert.Equal(t, uint64(1738754100), resp.Likers[0].UnixTimestamp)
		assert.Equal(t, "actor2", resp.Likers[1].ActorId)
		assert.Equal(t, uint64(1738686000), resp.Likers[1].UnixTimestamp)
		assert.Nil(t, resp.NextPaginationToken)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("Success_WithCursor", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		paginationToken := "eyJ1cGRhdGVkX2F0IjoiMjAyNS0wMi0wMiAxNTo0NTowMCIsImFjdG9yX2lkIjoiNyJ9"
		decodedCursor := &entity.Cursor{
			UpdatedAt: time.Date(2025, 2, 2, 15, 45, 0, 0, time.UTC),
			ActorId:   "7",
		}
		pageSize := 50
		likers := []entity.Liker{
			{ActorID: "actor3", UnixTimestamp: 1738400000},
			{ActorID: "actor4", UnixTimestamp: 1738300000},
		}
		nextCursor := &entity.Cursor{
			UpdatedAt: time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC),
			ActorId:   "actor4",
		}
		nextToken := "eyJ1cGRhdGVkX2F0IjoiMjAyNS0wMi0wMSAxMjowMDowMCIsImFjdG9yX2lkIjoiYWN0b3I0In0="

		// Set up mock expectations
		mockRepo.On("ListLikersByRecipient", ctx, recipientID, decodedCursor, pageSize).
			Return(likers, nextCursor, nil)
		mockCursorService.On("Decode", paginationToken).
			Return(decodedCursor, nil)
		mockCursorService.On("Encode", nextCursor).
			Return(nextToken, nil)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
			PaginationToken: &paginationToken,
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Likers, 2)
		assert.Equal(t, "actor3", resp.Likers[0].ActorId)
		assert.Equal(t, uint64(1738400000), resp.Likers[0].UnixTimestamp)
		assert.Equal(t, "actor4", resp.Likers[1].ActorId)
		assert.Equal(t, uint64(1738300000), resp.Likers[1].UnixTimestamp)
		assert.NotNil(t, resp.NextPaginationToken)
		assert.Equal(t, nextToken, *resp.NextPaginationToken)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("EmptyRecipientID", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request with empty recipient ID
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: "",
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "missing recipient user id")
		assert.Nil(t, resp)

		// Verify expectations (no calls expected)
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("InvalidPaginationToken", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		invalidToken := "invalid-token"

		// Set up mock expectations
		mockCursorService.On("Decode", invalidToken).
			Return(nil, errors.New("invalid token format"))

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request with invalid token
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
			PaginationToken: &invalidToken,
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "invalid pagination token")
		assert.Nil(t, resp)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		pageSize := 50

		// Set up mock expectations
		mockRepo.On("ListLikersByRecipient", ctx, recipientID, (*entity.Cursor)(nil), pageSize).
			Return(nil, nil, errors.New("database error"))

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to fetch likers")
		assert.Nil(t, resp)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("CursorEncodingError", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		pageSize := 50
		likers := []entity.Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
		}
		nextCursor := &entity.Cursor{
			UpdatedAt: time.Date(2025, 2, 1, 12, 0, 0, 0, time.UTC),
			ActorId:   "actor1",
		}

		// Set up mock expectations
		mockRepo.On("ListLikersByRecipient", ctx, recipientID, (*entity.Cursor)(nil), pageSize).
			Return(likers, nextCursor, nil)
		mockCursorService.On("Encode", nextCursor).
			Return("", errors.New("encoding failed"))

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.ListLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to encode pagination token")
		assert.Nil(t, resp)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})
}

func TestListNewLikedYou(t *testing.T) {
	// Create a context for testing
	ctx := context.Background()

	t.Run("Success_WithoutCursor", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		pageSize := 50
		likers := []entity.Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
			{ActorID: "actor2", UnixTimestamp: 1738686000},
		}

		// Set up mock expectations
		mockRepo.On("ListNewLikersByRecipient", ctx, recipientID, (*entity.Cursor)(nil), pageSize).
			Return(likers, nil, nil)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.ListLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.ListNewLikedYou(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Likers, 2)
		assert.Equal(t, "actor1", resp.Likers[0].ActorId)
		assert.Equal(t, uint64(1738754100), resp.Likers[0].UnixTimestamp)
		assert.Equal(t, "actor2", resp.Likers[1].ActorId)
		assert.Equal(t, uint64(1738686000), resp.Likers[1].UnixTimestamp)
		assert.Nil(t, resp.NextPaginationToken)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	// Additional test cases for ListNewLikedYou would be similar to ListLikedYou
	// We'll skip the duplicated tests for brevity
}

func TestCountLikedYou(t *testing.T) {
	// Create a context for testing
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"
		expectedCount := uint64(42)

		// Set up mock expectations
		mockRepo.On("CountLikersByRecipient", ctx, recipientID).
			Return(expectedCount, nil)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.CountLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.CountLikedYou(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, expectedCount, resp.Count)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("EmptyRecipientID", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request with empty recipient ID
		req := &grpclibs.CountLikedYouRequest{
			RecipientUserId: "",
		}

		// Call the method
		resp, err := server.CountLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "missing recipient user id")
		assert.Nil(t, resp)

		// Verify expectations (no calls expected)
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		recipientID := "user123"

		// Set up mock expectations
		mockRepo.On("CountLikersByRecipient", ctx, recipientID).
			Return(uint64(0), errors.New("database error"))

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.CountLikedYouRequest{
			RecipientUserId: recipientID,
		}

		// Call the method
		resp, err := server.CountLikedYou(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to count likers")
		assert.Nil(t, resp)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})
}

func TestPutDecision(t *testing.T) {
	// Create a context for testing
	ctx := context.Background()

	t.Run("Success_MutualLike", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		actorID := "actor123"
		recipientID := "recipient456"
		liked := true
		mutualLike := true

		// Set up mock expectations
		mockRepo.On("CreateOrUpdateDecision", ctx, actorID, recipientID, liked).
			Return(mutualLike, nil)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.PutDecisionRequest{
			ActorUserId:     actorID,
			RecipientUserId: recipientID,
			LikedRecipient:  liked,
		}

		// Call the method
		resp, err := server.PutDecision(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mutualLike, resp.MutualLikes)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("Success_NoMutualLike", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		actorID := "actor123"
		recipientID := "recipient456"
		liked := true
		mutualLike := false

		// Set up mock expectations
		mockRepo.On("CreateOrUpdateDecision", ctx, actorID, recipientID, liked).
			Return(mutualLike, nil)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.PutDecisionRequest{
			ActorUserId:     actorID,
			RecipientUserId: recipientID,
			LikedRecipient:  liked,
		}

		// Call the method
		resp, err := server.PutDecision(ctx, req)

		// Assert results
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mutualLike, resp.MutualLikes)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("EmptyActorID", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request with empty actor ID
		req := &grpclibs.PutDecisionRequest{
			ActorUserId:     "",
			RecipientUserId: "recipient456",
			LikedRecipient:  true,
		}

		// Call the method
		resp, err := server.PutDecision(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "missing actor user id")
		assert.Nil(t, resp)

		// Verify expectations (no calls expected)
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("EmptyRecipientID", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request with empty recipient ID
		req := &grpclibs.PutDecisionRequest{
			ActorUserId:     "actor123",
			RecipientUserId: "",
			LikedRecipient:  true,
		}

		// Call the method
		resp, err := server.PutDecision(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "missing recipient user id")
		assert.Nil(t, resp)

		// Verify expectations (no calls expected)
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Create mock repository and cursor service
		mockRepo := new(repository.MockDecisionRepository)
		mockCursorService := new(entity.MockCursorService)

		// Define test data
		actorID := "actor123"
		recipientID := "recipient456"
		liked := true

		// Set up mock expectations
		mockRepo.On("CreateOrUpdateDecision", ctx, actorID, recipientID, liked).
			Return(false, errors.New("database error"))

		// Create server with mocks
		server := NewExploreGRPCServer(mockRepo, mockCursorService)

		// Create request
		req := &grpclibs.PutDecisionRequest{
			ActorUserId:     actorID,
			RecipientUserId: recipientID,
			LikedRecipient:  liked,
		}

		// Call the method
		resp, err := server.PutDecision(ctx, req)

		// Assert results
		require.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to put decision")
		assert.Nil(t, resp)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockCursorService.AssertExpectations(t)
	})
}
