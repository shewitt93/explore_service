package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCursor_EncodeDecode(t *testing.T) {
	// Create test cases
	testCases := []struct {
		name     string
		cursor   *Cursor
		wantErr  bool
		nilInput bool
	}{
		{
			name: "Valid cursor",
			cursor: &Cursor{
				UpdatedAt: time.Date(2025, 2, 2, 15, 45, 0, 0, time.UTC),
				ActorId:   "user123",
			},
			wantErr:  false,
			nilInput: false,
		},
		{
			name:     "Nil cursor",
			cursor:   nil,
			wantErr:  true,
			nilInput: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip test if cursor is nil and we're expecting an error
			if tc.nilInput {
				// Test EncodeCursor directly
				token, err := EncodeCursor(tc.cursor)
				assert.Error(t, err)
				assert.Empty(t, token)

				// Test via CursorService
				service := NewCursorService()
				token, err = service.Encode(tc.cursor)
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			// Encode the cursor
			token, err := EncodeCursor(tc.cursor)

			// Check encoding results
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, token)

			// Decode and verify the cursor
			decodedCursor, err := DecodeCursor(token)
			require.NoError(t, err)

			// Compare original and decoded cursor
			assert.Equal(t, tc.cursor.UpdatedAt.UTC(), decodedCursor.UpdatedAt.UTC())
			assert.Equal(t, tc.cursor.ActorId, decodedCursor.ActorId)

			// Test using CursorService
			service := NewCursorService()

			// Encode using service
			serviceToken, err := service.Encode(tc.cursor)
			require.NoError(t, err)
			assert.Equal(t, token, serviceToken)

			// Decode using service
			serviceDecodedCursor, err := service.Decode(serviceToken)
			require.NoError(t, err)
			assert.Equal(t, tc.cursor.UpdatedAt.UTC(), serviceDecodedCursor.UpdatedAt.UTC())
			assert.Equal(t, tc.cursor.ActorId, serviceDecodedCursor.ActorId)
		})
	}
}

func TestDecodeCursor_InvalidInput(t *testing.T) {
	t.Run("Invalid base64", func(t *testing.T) {
		_, err := DecodeCursor("this-is-not-base64!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode cursor")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// Create invalid JSON in valid base64
		invalidJSON := "eyJpbnZhbGlkIjoianNvbn0=" // Base64 of {"invalid":"json"
		_, err := DecodeCursor(invalidJSON)
		assert.Error(t, err)
	})

	t.Run("Missing updated_at field", func(t *testing.T) {
		// Base64 of {"actor_id":"user123"}
		missingUpdatedAt := "eyJhY3Rvcl9pZCI6InVzZXIxMjMifQ=="
		_, err := DecodeCursor(missingUpdatedAt)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing or invalid 'updated_at' field")
	})

	t.Run("Missing actor_id field", func(t *testing.T) {
		// Base64 of {"updated_at":"2025-02-02 15:45:00"}
		missingActorId := "eyJ1cGRhdGVkX2F0IjoiMjAyNS0wMi0wMiAxNTo0NTowMCJ9"
		_, err := DecodeCursor(missingActorId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing or invalid 'actor_id' field")
	})

	t.Run("Invalid timestamp format", func(t *testing.T) {
		// Base64 of {"updated_at":"invalid-time","actor_id":"user123"}
		invalidTimestamp := "eyJ1cGRhdGVkX2F0IjoiaW52YWxpZC10aW1lIiwiYWN0b3JfaWQiOiJ1c2VyMTIzIn0="
		_, err := DecodeCursor(invalidTimestamp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid 'updated_at' format")
	})
}

func TestCreateNextCursor(t *testing.T) {
	t.Run("LikersLessThanLimit", func(t *testing.T) {
		// Prepare test data
		likers := []Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
			{ActorID: "actor2", UnixTimestamp: 1738686000},
		}
		limit := 5

		// Call function
		resultLikers, nextCursor := CreateNextCursor(likers, limit)

		// Verify results
		assert.Equal(t, likers, resultLikers)
		assert.Nil(t, nextCursor)

		// Test via service
		service := NewCursorService()
		resultLikers, nextCursor = service.CreateNextCursor(likers, limit)
		assert.Equal(t, likers, resultLikers)
		assert.Nil(t, nextCursor)
	})

	t.Run("LikersEqualToLimit", func(t *testing.T) {
		// Prepare test data
		likers := []Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
			{ActorID: "actor2", UnixTimestamp: 1738686000},
		}
		limit := 2

		// Call function
		resultLikers, nextCursor := CreateNextCursor(likers, limit)

		// Verify results
		assert.Equal(t, likers, resultLikers)
		assert.Nil(t, nextCursor)

		// Test via service
		service := NewCursorService()
		resultLikers, nextCursor = service.CreateNextCursor(likers, limit)
		assert.Equal(t, likers, resultLikers)
		assert.Nil(t, nextCursor)
	})

	t.Run("LikersMoreThanLimit", func(t *testing.T) {
		// Prepare test data
		likers := []Liker{
			{ActorID: "actor1", UnixTimestamp: 1738754100},
			{ActorID: "actor2", UnixTimestamp: 1738686000},
			{ActorID: "actor3", UnixTimestamp: 1738600000},
		}
		limit := 2

		// Call function
		resultLikers, nextCursor := CreateNextCursor(likers, limit)

		// Verify results
		assert.Len(t, resultLikers, 2)
		assert.Equal(t, likers[0], resultLikers[0])
		assert.Equal(t, likers[1], resultLikers[1])

		assert.NotNil(t, nextCursor)
		assert.Equal(t, time.Unix(1738686000, 0), nextCursor.UpdatedAt)
		assert.Equal(t, "actor2", nextCursor.ActorId)

		// Test via service
		service := NewCursorService()
		resultLikers, nextCursor = service.CreateNextCursor(likers, limit)
		assert.Len(t, resultLikers, 2)
		assert.Equal(t, likers[0], resultLikers[0])
		assert.Equal(t, likers[1], resultLikers[1])

		assert.NotNil(t, nextCursor)
		assert.Equal(t, time.Unix(1738686000, 0), nextCursor.UpdatedAt)
		assert.Equal(t, "actor2", nextCursor.ActorId)
	})

	t.Run("EmptyLikers", func(t *testing.T) {
		// Prepare test data
		var likers []Liker
		limit := 5

		// Call function
		resultLikers, nextCursor := CreateNextCursor(likers, limit)

		// Verify results
		assert.Empty(t, resultLikers)
		assert.Nil(t, nextCursor)

		// Test via service
		service := NewCursorService()
		resultLikers, nextCursor = service.CreateNextCursor(likers, limit)
		assert.Empty(t, resultLikers)
		assert.Nil(t, nextCursor)
	})
}

func TestDefaultCursorService(t *testing.T) {
	t.Run("NewCursorService", func(t *testing.T) {
		service := NewCursorService()
		assert.NotNil(t, service)
		assert.IsType(t, &DefaultCursorService{}, service)
	})
}
