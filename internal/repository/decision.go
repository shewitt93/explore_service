package repository

import (
	"context"
	"github.com/shewitt93/explore_service/internal/entity"
)

type DecisionRepository interface {
	ListLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error)

	ListNewLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error)

	CountLikersByRecipient(ctx context.Context, recipientID string) (uint64, error)

	CreateOrUpdateDecision(ctx context.Context, actorID string, recipientID string, liked bool) (bool, error)
}
