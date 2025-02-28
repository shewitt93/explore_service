package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shewitt93/explore_service/internal/entity"
)

type DecisionRepositoryImpl struct {
	db *sql.DB
}

func NewDecisionRepositoryImpl(db *sql.DB) DecisionRepository {
	return DecisionRepositoryImpl{
		db: db,
	}
}

func (r DecisionRepositoryImpl) ListLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error) {

	baseQuery := "SELECT actor_id, UNIX_TIMESTAMP(updated_at) as unix_timestamp FROM user_decisions WHERE recipient_id = ? AND liked = TRUE"

	var args []interface{}
	var query string

	// Apply cursor pagination if provided
	if cursor != nil {
		query = baseQuery + `
			AND (updated_at < ? OR (updated_at = ? AND actor_id < ?))
			ORDER BY updated_at DESC, actor_id DESC
			LIMIT ?`
		args = []interface{}{recipientID, cursor.UpdatedAt, cursor.UpdatedAt, cursor.ActorId, limit + 1}
	} else {
		query = baseQuery + `
			ORDER BY updated_at DESC, actor_id DESC
			LIMIT ?`
		args = []interface{}{recipientID, limit + 1}
	}

	// Execute the query and get results
	likers, err := r.executeLikersQuery(ctx, query, args)
	if err != nil {
		return nil, nil, err
	}

	// Handle pagination only if a cursor was originally provided
	var nextCursor *entity.Cursor
	if cursor != nil && len(likers) > limit {
		likers, nextCursor = entity.CreateNextCursor(likers, limit)
	} else if len(likers) > limit {
		// Just trim the results without creating a cursor
		likers = likers[:limit]
	}

	return likers, nextCursor, nil
}

func (r DecisionRepositoryImpl) ListNewLikersByRecipient(ctx context.Context, recipientID string, cursor *entity.Cursor, limit int) ([]entity.Liker, *entity.Cursor, error) {
	// Base query excluding mutual likes
	query := `
        SELECT d1.actor_id, UNIX_TIMESTAMP(d1.updated_at) as unix_timestamp
        FROM user_decisions d1
        LEFT JOIN user_decisions d2 
            ON d1.actor_id = d2.recipient_id 
            AND d2.actor_id = d1.recipient_id 
            AND d2.liked = TRUE
        WHERE d1.recipient_id = ? AND d1.liked = TRUE AND d2.actor_id IS NULL`
	args := []interface{}{recipientID}

	// Apply cursor-based pagination if provided
	if cursor != nil {
		query += ` AND (d1.updated_at < ? OR (d1.updated_at = ? AND d1.actor_id < ?))`
		args = append(args, cursor.UpdatedAt, cursor.UpdatedAt, cursor.ActorId)
	}

	query += ` ORDER BY d1.updated_at DESC, d1.actor_id DESC LIMIT ?`
	args = append(args, limit+1) // Fetch one extra to check for next page

	// Execute the query and get results
	likers, err := r.executeLikersQuery(ctx, query, args)
	if err != nil {
		return nil, nil, err
	}

	// Handle pagination only if a cursor was originally provided
	var nextCursor *entity.Cursor
	if cursor != nil && len(likers) > limit {
		likers, nextCursor = entity.CreateNextCursor(likers, limit)
	} else if len(likers) > limit {
		// Just trim the results without creating a cursor
		likers = likers[:limit]
	}

	return likers, nextCursor, nil
}

// executeLikersQuery executes the SQL query and transforms the results into entities
func (r DecisionRepositoryImpl) executeLikersQuery(ctx context.Context, query string, args []interface{}) ([]entity.Liker, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var likers []entity.Liker
	for rows.Next() {
		var liker entity.Liker
		var unixTs int64
		if err := rows.Scan(&liker.ActorID, &unixTs); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		liker.UnixTimestamp = uint64(unixTs)
		likers = append(likers, liker)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return likers, nil
}

func (r DecisionRepositoryImpl) CountLikersByRecipient(ctx context.Context, recipientID string) (uint64, error) {
	query := "SELECT COUNT(*) FROM user_decisions WHERE recipient_id = ? AND liked = TRUE"

	var count uint64
	err := r.db.QueryRowContext(ctx, query, recipientID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count likers: %w", err)
	}

	return count, nil
}

func (r DecisionRepositoryImpl) CreateOrUpdateDecision(ctx context.Context, actorID string, recipientID string, liked bool) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Insert or update the decision
	query := `
		INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE liked = ?, updated_at = NOW()`

	_, err = tx.ExecContext(ctx, query, actorID, recipientID, liked, liked)
	if err != nil {
		return false, fmt.Errorf("failed to put decision: %w", err)
	}

	// If the decision is a like, check if there's a mutual like
	mutualLike := false
	if liked {
		checkQuery := `
			SELECT EXISTS(
				SELECT 1
				FROM user_decisions
				WHERE actor_id = ? AND recipient_id = ? AND liked = TRUE
			)`

		err = tx.QueryRowContext(ctx, checkQuery, recipientID, actorID).Scan(&mutualLike)
		if err != nil {
			return false, fmt.Errorf("failed to check for mutual like: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mutualLike, nil
}
