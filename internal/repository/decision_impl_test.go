package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shewitt93/explore_service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLikersByRecipient(t *testing.T) {
	// Create a new mock database connection with QueryMatcherEqual
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock db
	repo := NewDecisionRepositoryImpl(db)

	// Create a test context
	ctx := context.Background()

	t.Run("Success_WithoutCursor", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		limit := 10

		// Setup expected query and response
		rows := sqlmock.NewRows([]string{"actor_id", "unix_timestamp"}).
			AddRow("actor1", int64(1738754100)).
			AddRow("actor2", int64(1738686000))

		// Define expected SQL with args
		expectedSQL := "SELECT actor_id, UNIX_TIMESTAMP(updated_at) as unix_timestamp FROM user_decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY updated_at DESC, actor_id DESC LIMIT ?"
		mock.ExpectQuery(expectedSQL).
			WithArgs(recipientID, limit+1).
			WillReturnRows(rows)

		// Call the method
		likers, nextCursor, err := repo.ListLikersByRecipient(ctx, recipientID, nil, limit)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, likers, 2)
		assert.Equal(t, "actor1", likers[0].ActorID)
		assert.Equal(t, uint64(1738754100), likers[0].UnixTimestamp)
		assert.Equal(t, "actor2", likers[1].ActorID)
		assert.Equal(t, uint64(1738686000), likers[1].UnixTimestamp)
		assert.Nil(t, nextCursor)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Success_WithCursor", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		limit := 2
		cursorTime := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
		cursor := &entity.Cursor{
			UpdatedAt: cursorTime,
			ActorId:   "actor3",
		}

		// Setup expected query and response
		rows := sqlmock.NewRows([]string{"actor_id", "unix_timestamp"}).
			AddRow("actor4", int64(1738400000)).
			AddRow("actor5", int64(1738300000)).
			AddRow("actor6", int64(1738200000)) // Extra row to test pagination

		// Define expected SQL with args
		expectedSQL := "SELECT actor_id, UNIX_TIMESTAMP(updated_at) as unix_timestamp FROM user_decisions WHERE recipient_id = ? AND liked = TRUE AND (updated_at < ? OR (updated_at = ? AND actor_id < ?)) ORDER BY updated_at DESC, actor_id DESC LIMIT ?"
		mock.ExpectQuery(expectedSQL).
			WithArgs(recipientID, cursorTime, cursorTime, "actor3", limit+1).
			WillReturnRows(rows)

		// Call the method
		likers, nextCursor, err := repo.ListLikersByRecipient(ctx, recipientID, cursor, limit)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, likers, 2)
		assert.Equal(t, "actor4", likers[0].ActorID)
		assert.Equal(t, "actor5", likers[1].ActorID)
		assert.NotNil(t, nextCursor)
		assert.Equal(t, time.Unix(1738300000, 0), nextCursor.UpdatedAt)
		assert.Equal(t, "actor5", nextCursor.ActorId)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		limit := 10

		// Setup expected query to return an error
		expectedSQL := "SELECT actor_id, UNIX_TIMESTAMP(updated_at) as unix_timestamp FROM user_decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY updated_at DESC, actor_id DESC LIMIT ?"
		mock.ExpectQuery(expectedSQL).
			WithArgs(recipientID, limit+1).
			WillReturnError(errors.New("database error"))

		// Call the method
		likers, nextCursor, err := repo.ListLikersByRecipient(ctx, recipientID, nil, limit)

		// Assert results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "database query failed")
		assert.Nil(t, likers)
		assert.Nil(t, nextCursor)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestListNewLikersByRecipient(t *testing.T) {
	// Create a test context
	ctx := context.Background()

	t.Run("Success_WithoutCursor", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		limit := 10

		// Create mock database connection with default regexp matcher
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewDecisionRepositoryImpl(db)

		// Setup expected query and response
		rows := sqlmock.NewRows([]string{"actor_id", "unix_timestamp"}).
			AddRow("actor1", int64(1738754100)).
			AddRow("actor2", int64(1738686000))

		// Since this query is complex with whitespace variations, we have to use substring matching
		mock.ExpectQuery("SELECT d1.actor_id, UNIX_TIMESTAMP").
			WithArgs(recipientID, limit+1).
			WillReturnRows(rows)

		// Call the method
		likers, nextCursor, err := repo.ListNewLikersByRecipient(ctx, recipientID, nil, limit)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, likers, 2)
		assert.Equal(t, "actor1", likers[0].ActorID)
		assert.Equal(t, uint64(1738754100), likers[0].UnixTimestamp)
		assert.Equal(t, "actor2", likers[1].ActorID)
		assert.Equal(t, uint64(1738686000), likers[1].UnixTimestamp)
		assert.Nil(t, nextCursor)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Success_WithCursor", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		limit := 2
		cursorTime := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)
		cursor := &entity.Cursor{
			UpdatedAt: cursorTime,
			ActorId:   "actor3",
		}

		// Create mock database connection with default regexp matcher
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewDecisionRepositoryImpl(db)

		// Setup expected query and response
		rows := sqlmock.NewRows([]string{"actor_id", "unix_timestamp"}).
			AddRow("actor4", int64(1738400000)).
			AddRow("actor5", int64(1738300000)).
			AddRow("actor6", int64(1738200000)) // Extra row to test pagination

		// Since this query is complex with whitespace variations, we have to use substring matching
		mock.ExpectQuery("SELECT d1.actor_id, UNIX_TIMESTAMP").
			WithArgs(recipientID, cursorTime, cursorTime, "actor3", limit+1).
			WillReturnRows(rows)

		// Call the method
		likers, nextCursor, err := repo.ListNewLikersByRecipient(ctx, recipientID, cursor, limit)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, likers, 2)
		assert.Equal(t, "actor4", likers[0].ActorID)
		assert.Equal(t, "actor5", likers[1].ActorID)
		assert.NotNil(t, nextCursor)
		assert.Equal(t, time.Unix(1738300000, 0), nextCursor.UpdatedAt)
		assert.Equal(t, "actor5", nextCursor.ActorId)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCountLikersByRecipient(t *testing.T) {
	// Create a new mock database connection with QueryMatcherEqual
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock db
	repo := NewDecisionRepositoryImpl(db)

	// Create a test context
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"
		expectedCount := uint64(5)

		// Setup expected query and response
		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(expectedCount)

		// Define expected SQL with args
		expectedSQL := "SELECT COUNT(*) FROM user_decisions WHERE recipient_id = ? AND liked = TRUE"
		mock.ExpectQuery(expectedSQL).
			WithArgs(recipientID).
			WillReturnRows(rows)

		// Call the method
		count, err := repo.CountLikersByRecipient(ctx, recipientID)

		// Assert results
		require.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Define test data
		recipientID := "recipient1"

		// Setup expected query to return an error
		expectedSQL := "SELECT COUNT(*) FROM user_decisions WHERE recipient_id = ? AND liked = TRUE"
		mock.ExpectQuery(expectedSQL).
			WithArgs(recipientID).
			WillReturnError(errors.New("database error"))

		// Call the method
		count, err := repo.CountLikersByRecipient(ctx, recipientID)

		// Assert results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to count likers")
		assert.Equal(t, uint64(0), count)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestCreateOrUpdateDecision(t *testing.T) {
	// Create a new mock database connection with QueryMatcherEqual
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock db
	repo := NewDecisionRepositoryImpl(db)

	// Create a test context
	ctx := context.Background()

	t.Run("Success_MutualLike", func(t *testing.T) {
		// Define test data
		actorID := "actor1"
		recipientID := "recipient1"
		liked := true

		// Setup transaction expectations
		mock.ExpectBegin()

		// Setup insert/update expectation
		expectedSQL := "INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE liked = ?, updated_at = NOW()"
		mock.ExpectExec(expectedSQL).
			WithArgs(actorID, recipientID, liked, liked).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup mutual like check expectation
		checkSQL := "SELECT EXISTS( SELECT 1 FROM user_decisions WHERE actor_id = ? AND recipient_id = ? AND liked = TRUE )"
		checkRows := sqlmock.NewRows([]string{"exists"}).AddRow(1) // 1 means mutual like exists
		mock.ExpectQuery(checkSQL).
			WithArgs(recipientID, actorID).
			WillReturnRows(checkRows)

		// Setup commit expectation
		mock.ExpectCommit()

		// Call the method
		mutualLike, err := repo.CreateOrUpdateDecision(ctx, actorID, recipientID, liked)

		// Assert results
		require.NoError(t, err)
		assert.True(t, mutualLike)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Success_NoMutualLike", func(t *testing.T) {
		// Define test data
		actorID := "actor1"
		recipientID := "recipient2"
		liked := true

		// Setup transaction expectations
		mock.ExpectBegin()

		// Setup insert/update expectation
		expectedSQL := "INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE liked = ?, updated_at = NOW()"
		mock.ExpectExec(expectedSQL).
			WithArgs(actorID, recipientID, liked, liked).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup mutual like check expectation
		checkSQL := "SELECT EXISTS( SELECT 1 FROM user_decisions WHERE actor_id = ? AND recipient_id = ? AND liked = TRUE )"
		checkRows := sqlmock.NewRows([]string{"exists"}).AddRow(0) // 0 means no mutual like
		mock.ExpectQuery(checkSQL).
			WithArgs(recipientID, actorID).
			WillReturnRows(checkRows)

		// Setup commit expectation
		mock.ExpectCommit()

		// Call the method
		mutualLike, err := repo.CreateOrUpdateDecision(ctx, actorID, recipientID, liked)

		// Assert results
		require.NoError(t, err)
		assert.False(t, mutualLike)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Success_Pass", func(t *testing.T) {
		// Define test data
		actorID := "actor1"
		recipientID := "recipient3"
		liked := false // This is a pass, not a like

		// Setup transaction expectations
		mock.ExpectBegin()

		// Setup insert/update expectation
		expectedSQL := "INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE liked = ?, updated_at = NOW()"
		mock.ExpectExec(expectedSQL).
			WithArgs(actorID, recipientID, liked, liked).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// No mutual like check for pass decisions

		// Setup commit expectation
		mock.ExpectCommit()

		// Call the method
		mutualLike, err := repo.CreateOrUpdateDecision(ctx, actorID, recipientID, liked)

		// Assert results
		require.NoError(t, err)
		assert.False(t, mutualLike) // A pass can't create a mutual like

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("TransactionError", func(t *testing.T) {
		// Define test data
		actorID := "actor1"
		recipientID := "recipient4"
		liked := true

		// Setup transaction to fail
		mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

		// Call the method
		mutualLike, err := repo.CreateOrUpdateDecision(ctx, actorID, recipientID, liked)

		// Assert results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to begin transaction")
		assert.False(t, mutualLike)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		// Define test data
		actorID := "actor1"
		recipientID := "recipient5"
		liked := true

		// Setup transaction expectations
		mock.ExpectBegin()

		// Setup insert/update to fail
		expectedSQL := "INSERT INTO user_decisions (actor_id, recipient_id, liked, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE liked = ?, updated_at = NOW()"
		mock.ExpectExec(expectedSQL).
			WithArgs(actorID, recipientID, liked, liked).
			WillReturnError(errors.New("query error"))

		// Setup rollback expectation
		mock.ExpectRollback()

		// Call the method
		mutualLike, err := repo.CreateOrUpdateDecision(ctx, actorID, recipientID, liked)

		// Assert results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to put decision")
		assert.False(t, mutualLike)

		// Ensure all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestExecuteLikersQuery_ScanError(t *testing.T) {
	// Create a new mock database connection with QueryMatcherEqual
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	// Create the repository with the mock db
	repo := NewDecisionRepositoryImpl(db)

	// Create a test context
	ctx := context.Background()

	// Define test data
	recipientID := "recipient1"
	limit := 10

	// Setup expected query with type mismatch to force scan error
	rows := sqlmock.NewRows([]string{"actor_id", "unix_timestamp"}).
		AddRow("actor1", "not a number") // This will cause a scan error

	// Define expected SQL with args
	expectedSQL := "SELECT actor_id, UNIX_TIMESTAMP(updated_at) as unix_timestamp FROM user_decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY updated_at DESC, actor_id DESC LIMIT ?"
	mock.ExpectQuery(expectedSQL).
		WithArgs(recipientID, limit+1).
		WillReturnRows(rows)

	// Call the method
	likers, nextCursor, err := repo.ListLikersByRecipient(ctx, recipientID, nil, limit)

	// Assert results
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan row")
	assert.Nil(t, likers)
	assert.Nil(t, nextCursor)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
