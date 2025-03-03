package repository

import (
	"context"

	"github.com/shewitt93/explore_service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockDecisionRepository is a mock implementation of DecisionRepository
type MockDecisionRepository struct {
	mock.Mock
}

// Ensure MockDecisionRepository implements DecisionRepository interface
var _ DecisionRepository = (*MockDecisionRepository)(nil)

func (m *MockDecisionRepository) ListLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error) {
	args := m.Called(ctx, recipientID, cursor, limit)

	var likers []entity.Liker
	if args.Get(0) != nil {
		likers = args.Get(0).([]entity.Liker)
	}

	var nextCursor *entity.Cursor
	if args.Get(1) != nil {
		nextCursor = args.Get(1).(*entity.Cursor)
	}

	return likers, nextCursor, args.Error(2)
}

func (m *MockDecisionRepository) ListNewLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error) {
	args := m.Called(ctx, recipientID, cursor, limit)

	var likers []entity.Liker
	if args.Get(0) != nil {
		likers = args.Get(0).([]entity.Liker)
	}

	var nextCursor *entity.Cursor
	if args.Get(1) != nil {
		nextCursor = args.Get(1).(*entity.Cursor)
	}

	return likers, nextCursor, args.Error(2)
}

func (m *MockDecisionRepository) CountLikersByRecipient(ctx context.Context, recipientID string) (uint64, error) {
	args := m.Called(ctx, recipientID)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDecisionRepository) CreateOrUpdateDecision(ctx context.Context, actorID string, recipientID string, liked bool) (bool, error) {
	args := m.Called(ctx, actorID, recipientID, liked)
	return args.Bool(0), args.Error(1)
}
