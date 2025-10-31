package lightforge

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pballentine13/liteforge"
)

// setupTestDB initializes an in-memory SQLite database, creates the TestUser table,
// and returns a liteforge.Repository and a cleanup function.
func setupTestDB(t *testing.T) (liteforge.Repository, func()) {
	cfg := liteforge.Config{
		DriverName:     "sqlite3",
		DataSourceName: ":memory:",
	}

	ds, err := liteforge.OpenDB(cfg)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the table
	err = liteforge.CreateTable(ds, TestUser{})
	if err != nil {
		ds.DB.Close()
		t.Fatalf("Failed to create table: %v", err)
	}

	repo := liteforge.NewRepository(ds)

	cleanup := func() {
		ds.DB.Close()
	}

	return repo, cleanup
}

func TestRepository_Save_Insert(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	user := &TestUser{
		Username: "testuser_insert",
		Email:    "insert@example.com",
		Age:      30,
		IsActive: true,
	}

	// 1. Test INSERT (ID is zero)
	result, err := repo.Save(user)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	lastID, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.True(t, lastID > 0)

	// 2. Verify insertion by finding the record
	foundUser := &TestUser{}
	err = repo.FindByID(foundUser, int(lastID))
	assert.NoError(t, err)
	assert.Equal(t, int(lastID), foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Age, foundUser.Age)
	assert.Equal(t, user.IsActive, foundUser.IsActive)
}

func TestRepository_Save_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup: Insert a record first
	user := &TestUser{
		Username: "testuser_update",
		Email:    "original@example.com",
		Age:      25,
		IsActive: false,
	}
	result, err := repo.Save(user)
	assert.NoError(t, err)
	lastID, _ := result.LastInsertId()

	// Manually set the ID on the model for the update operation
	user.ID = int(lastID)

	// 1. Test UPDATE (ID is non-zero)
	user.Email = "updated@example.com"
	user.Age = 50
	user.IsActive = true

	result, err = repo.Save(user) // Should call Update because ID is non-zero
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// 2. Verify update by finding the record
	foundUser := &TestUser{}
	err = repo.FindByID(foundUser, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "updated@example.com", foundUser.Email)
	assert.Equal(t, 50, foundUser.Age)
	assert.Equal(t, true, foundUser.IsActive)
}

func TestRepository_FindByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup: Insert a record
	user := &TestUser{
		Username: "find_test",
		Email:    "find@example.com",
		Age:      40,
		IsActive: true,
	}
	result, _ := repo.Save(user)
	lastID, _ := result.LastInsertId()

	// 1. Test successful retrieval
	t.Run("Success", func(t *testing.T) {
		foundUser := &TestUser{}
		err := repo.FindByID(foundUser, int(lastID))
		assert.NoError(t, err)
		assert.Equal(t, int(lastID), foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
	})

	// 2. Test record not found
	t.Run("NotFound", func(t *testing.T) {
		notFoundUser := &TestUser{}
		// Use an ID that definitely doesn't exist
		err := repo.FindByID(notFoundUser, 99999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		// Ensure the model is not populated (ID should be zero)
		assert.Equal(t, 0, notFoundUser.ID)
	})

	// 3. Test invalid model input (edge case/error path)
	t.Run("InvalidModel", func(t *testing.T) {
		err := repo.FindByID(TestUser{}, int(lastID)) // Pass struct instead of pointer
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "model must be a non-nil pointer to a struct")
	})
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup: Insert a record
	user := &TestUser{
		Username: "update_explicit",
		Email:    "pre_update@example.com",
		Age:      10,
		IsActive: false,
	}
	result, _ := repo.Save(user)
	lastID, _ := result.LastInsertId()
	user.ID = int(lastID)

	// 1. Perform explicit update
	user.Email = "post_update@example.com"
	user.Age = 100

	result, err := repo.Update(user)
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// 2. Verify the update
	foundUser := &TestUser{}
	err = repo.FindByID(foundUser, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "post_update@example.com", foundUser.Email)
	assert.Equal(t, 100, foundUser.Age)

	// 3. Test update on non-existent ID (should affect 0 rows)
	t.Run("NonExistentID", func(t *testing.T) {
		nonExistentUser := &TestUser{ID: 99999, Username: "foo", Email: "bar@baz.com"}
		result, err := repo.Update(nonExistentUser)
		assert.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)
	})
}

func TestRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup: Insert a record
	user := &TestUser{
		Username: "delete_test",
		Email:    "delete@example.com",
		Age:      60,
		IsActive: true,
	}
	result, _ := repo.Save(user)
	lastID, _ := result.LastInsertId()
	user.ID = int(lastID)

	// 1. Test successful deletion
	result, err := repo.Delete(user)
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// 2. Verify deletion by trying to find the record
	foundUser := &TestUser{}
	err = repo.FindByID(foundUser, user.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	// 3. Test deleting a non-existent record (should affect 0 rows)
	t.Run("NonExistentID", func(t *testing.T) {
		nonExistentUser := &TestUser{ID: 99999}
		result, err := repo.Delete(nonExistentUser)
		assert.NoError(t, err)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(0), rowsAffected)
	})
}
