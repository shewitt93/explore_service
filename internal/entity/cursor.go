package entity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Cursor struct {
	UpdatedAt time.Time
	ActorId   string
}

func DecodeCursor(cursor string) (*Cursor, error) {
	// Decode base64-encoded cursor string
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	// Unmarshal JSON into a map
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor data: %w", err)
	}

	// Extract "updated_at" field
	updatedAtStr, ok := dataMap["updated_at"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'updated_at' field")
	}

	// Parse the timestamp
	updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'updated_at' format: %w", err)
	}

	// Extract "actor_id" field
	actorId, ok := dataMap["actor_id"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'actor_id' field")
	}

	// Construct and return Cursor object
	return &Cursor{
		UpdatedAt: updatedAt,
		ActorId:   actorId,
	}, nil
}

func EncodeCursor(c *Cursor) (string, error) {
	if c == nil {
		return "", errors.New("cursor cannot be nil")
	}

	// Create a map for the cursor data
	dataMap := map[string]interface{}{
		"updated_at": c.UpdatedAt.Format("2006-01-02 15:04:05"),
		"actor_id":   c.ActorId,
	}

	// Marshal to JSON
	data, err := json.Marshal(dataMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor data: %w", err)
	}

	// Encode as base64
	return base64.StdEncoding.EncodeToString(data), nil
}

func CreateNextCursor(likers []Liker, limit int) ([]Liker, *Cursor) {
	// If we have more results than the limit, we need to create a cursor
	if len(likers) > limit {
		lastLiker := likers[limit-1]

		// Create cursor for the next page
		nextCursor := &Cursor{
			UpdatedAt: time.Unix(int64(lastLiker.UnixTimestamp), 0),
			ActorId:   lastLiker.ActorID,
		}

		// Trim results to the requested limit
		return likers[:limit], nextCursor
	}

	// No need for a cursor if we have fewer results than the limit
	return likers, nil
}
