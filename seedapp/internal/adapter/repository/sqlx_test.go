package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockDB creates a new mock database for testing
func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return sqlxDB, mock
}

func TestSqlxDB_Ping(t *testing.T) {
	// Setup
	mockDB, _ := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test case 1: Successful ping
	// Note: We can't directly test Ping with sqlmock as it doesn't support monitoring pings
	// Instead, we'll just verify that our function calls through to the underlying DB
	err := sqlxRepo.Ping()
	assert.NoError(t, err)

	// Test case 2: Failed ping
	// We can't easily simulate a ping error with sqlmock
	// In a real test environment, you might use a real database connection
	// or a more sophisticated mock
}

func TestSqlxDB_GetDB(t *testing.T) {
	// Setup
	mockDB, _ := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test
	db := sqlxRepo.GetDB()
	assert.Equal(t, mockDB, db)
}

func TestSqlxDB_SelectContext(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test data
	type TestStruct struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	testData := []TestStruct{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	query := "SELECT id, name FROM test_table"
	ctx := context.Background()

	// Test case: Successful select
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, d := range testData {
		rows.AddRow(d.ID, d.Name)
	}
	mock.ExpectQuery("SELECT id, name FROM test_table").WillReturnRows(rows)

	var result []TestStruct
	err := sqlxRepo.SelectContext(ctx, &result, query)
	assert.NoError(t, err)
	assert.Equal(t, testData, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_GetContext(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test data
	type TestStruct struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	testData := TestStruct{ID: 1, Name: "Test 1"}
	query := "SELECT id, name FROM test_table WHERE id = $1"
	ctx := context.Background()

	// Test case: Successful get
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(testData.ID, testData.Name)
	mock.ExpectQuery("SELECT id, name FROM test_table WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	var result TestStruct
	err := sqlxRepo.GetContext(ctx, &result, query, 1)
	assert.NoError(t, err)
	assert.Equal(t, testData, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_ExecContext(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test data
	query := "INSERT INTO test_table (name) VALUES ($1)"
	ctx := context.Background()
	expectedResult := sqlmock.NewResult(1, 1)

	// Test case 1: Without transaction
	mock.ExpectExec("INSERT INTO test_table \\(name\\) VALUES \\(\\$1\\)").
		WithArgs("test").
		WillReturnResult(expectedResult)

	result, err := sqlxRepo.ExecContext(ctx, query, "test")
	assert.NoError(t, err)
	rowsAffected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test case 2: With transaction
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test_table \\(name\\) VALUES \\(\\$1\\)").
		WithArgs("test").
		WillReturnResult(expectedResult)
	mock.ExpectCommit()

	tx, err := mockDB.Beginx()
	assert.NoError(t, err)
	txCtx := context.WithValue(ctx, TX_KEY, tx)

	result, err = sqlxRepo.ExecContext(txCtx, query, "test")
	assert.NoError(t, err)
	rowsAffected, _ = result.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected)

	err = tx.Commit()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_QueryRowxContext(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test data
	query := "SELECT id, name FROM test_table WHERE id = $1"
	ctx := context.Background()

	// Test case 1: Without transaction
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test")
	mock.ExpectQuery("SELECT id, name FROM test_table WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	row := sqlxRepo.QueryRowxContext(ctx, query, 1)
	assert.NotNil(t, row)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test case 2: With transaction
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test")
	mock.ExpectQuery("SELECT id, name FROM test_table WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := mockDB.Beginx()
	assert.NoError(t, err)
	txCtx := context.WithValue(ctx, TX_KEY, tx)

	row = sqlxRepo.QueryRowxContext(txCtx, query, 1)
	assert.NotNil(t, row)

	err = tx.Commit()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_Tx(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	ctx := context.Background()

	// Test case 1: Successful transaction
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test_table").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := sqlxRepo.Tx(ctx, func(ctx context.Context, errChan chan error) {
		_, err := sqlxRepo.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')")
		errChan <- err
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test case 2: Failed transaction
	mock.ExpectBegin()
	expectedErr := errors.New("exec error")
	mock.ExpectExec("INSERT INTO test_table").WillReturnError(expectedErr)
	mock.ExpectRollback()

	err = sqlxRepo.Tx(ctx, func(ctx context.Context, errChan chan error) {
		_, err := sqlxRepo.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')")
		errChan <- err
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_TxWithValue(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	ctx := context.Background()

	// Test case 1: Successful transaction with value
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test_table").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := sqlxRepo.TxWithValue(ctx, func(ctx context.Context, errChan chan error, resChan chan interface{}) {
		res, err := sqlxRepo.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')")
		errChan <- err
		resChan <- res
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test case 2: Failed transaction with value
	mock.ExpectBegin()
	expectedErr := errors.New("exec error")
	mock.ExpectExec("INSERT INTO test_table").WillReturnError(expectedErr)
	mock.ExpectRollback()

	result, err = sqlxRepo.TxWithValue(ctx, func(ctx context.Context, errChan chan error, resChan chan interface{}) {
		res, err := sqlxRepo.ExecContext(ctx, "INSERT INTO test_table (name) VALUES ('test')")
		errChan <- err
		resChan <- res
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlxDB_IsErrorDuplicate(t *testing.T) {
	// Setup
	mockDB, _ := setupMockDB(t)
	defer mockDB.Close()

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Test cases
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-duplicate error",
			err:      errors.New("some error"),
			expected: false,
		},
		// Note: We can't easily test the true case without importing pgx/pgconn
		// and creating a real PgError, which would add complexity to the test.
		// In a real implementation, you might want to add a test with a mock PgError.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sqlxRepo.IsErrorDuplicate(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSqlxDB_Startup(t *testing.T) {
	// This test is more complex as it involves actual connection to a database
	// In a real test, you might want to use a test database or mock the connection

	// For demonstration, we'll create a simple test with a mock config
	// We're skipping this test, so we don't need to create the sqlxRepo variable
	t.Skip("Skipping test that requires a real database connection")
}

func TestSqlxDB_Shutdown(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)

	sqlxRepo := &SqlxDB{
		DB: mockDB,
	}

	// Expect the Close call
	mock.ExpectClose()

	// Test
	err := sqlxRepo.Shutdown()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
