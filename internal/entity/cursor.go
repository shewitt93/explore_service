package entity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Cursor struct {
	UpdatedAt       time.Time
	RecipientUserId string
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
	updatedAt, err := time.Parse("2006-01-02 12:00:00", updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid 'updated_at' format: %w", err)
	}

	// Extract "recipient_user_id" field
	recipientUserId, ok := dataMap["recipient_user_id"].(string)
	if !ok {
		return nil, errors.New("missing or invalid 'recipient_user_id' field")
	}

	// Construct and return Cursor object
	return &Cursor{
		UpdatedAt:       updatedAt,
		RecipientUserId: recipientUserId,
	}, nil
}
