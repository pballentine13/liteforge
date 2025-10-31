package model

import (
	"database/sql"
	"errors"
	"testing"
)

// --- Mock Implementations for ORMDataStore Testing ---

// MockResult implements sql.Result
type MockResult struct {
	LastInsertID      int64
	RowsAffectedValue int64
}

func (m MockResult) LastInsertId() (int64, error) {
	return m.LastInsertID, nil
}

func (m MockResult) RowsAffected() (int64, error) {
	return m.RowsAffectedValue, nil
}

// MockRepository is a mock implementation of the Repository interface for testing.
type MockRepository struct {
	FindByIDFn func(model any, id int) error
	SaveFn     func(model any) (sql.Result, error)
	DeleteFn   func(model any) (sql.Result, error)
	UpdateFn   func(model any) (sql.Result, error)
}

func (m *MockRepository) FindByID(model any, id int) error {
	return m.FindByIDFn(model, id)
}
func (m *MockRepository) Save(model any) (sql.Result, error) {
	return m.SaveFn(model)
}
func (m *MockRepository) Delete(model any) (sql.Result, error) {
	return m.DeleteFn(model)
}
func (m *MockRepository) Update(model any) (sql.Result, error) {
	return m.UpdateFn(model)
}

// --- Tests for ORMDataStore ---

func TestORMDataStore_GetUserByID(t *testing.T) {
	t.Parallel()

	// Test Case 1: Successful retrieval
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockRepo := &MockRepository{
			FindByIDFn: func(model any, id int) error {
				user, ok := model.(*User) // Changed to *User
				if !ok {
					t.Fatalf("expected *User, got %T", model)
				}
				// Simulate populating the model
				user.ID = id
				user.Name = "Test User"
				user.Age = 42
				return nil
			},
		}
		ds := NewORMDataStore(mockRepo) // Changed to NewORMDataStore

		user, err := ds.GetUserByID(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.ID != 1 || user.Name != "Test User" {
			t.Errorf("unexpected user data: %+v", user)
		}
	})

	// Test Case 2: Not Found (sql.ErrNoRows)
	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		mockRepo := &MockRepository{
			FindByIDFn: func(model any, id int) error {
				return sql.ErrNoRows
			},
		}
		ds := NewORMDataStore(mockRepo)

		user, err := ds.GetUserByID(2)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if user != nil {
			t.Fatalf("expected nil user, got %+v", user)
		}
	})

	// Test Case 3: Underlying Repository Error
	t.Run("RepoError", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("database connection failed")
		mockRepo := &MockRepository{
			FindByIDFn: func(model any, id int) error {
				return expectedErr
			},
		}
		ds := NewORMDataStore(mockRepo)

		user, err := ds.GetUserByID(3)

		if user != nil {
			t.Fatalf("expected nil user, got %+v", user)
		}
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error to wrap %v, got %v", expectedErr, err)
		}
	})

	// Test Case 4: Nil Repository
	t.Run("NilRepo", func(t *testing.T) {
		t.Parallel()
		ds := &ORMDataStore{Repo: nil} // Changed to ORMDataStore

		_, err := ds.GetUserByID(4)

		if err == nil {
			t.Fatal("expected an error for nil repository, got nil")
		}
		if err.Error() != "repository is nil" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestORMDataStore_SaveUser(t *testing.T) {
	t.Parallel()

	// Test Case 1: Successful Save (Insert or Update)
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userToSave := &User{ID: 10, Name: "New Name", Age: 30} // Changed to *User

		mockRepo := &MockRepository{
			SaveFn: func(model any) (sql.Result, error) {
				savedUser, ok := model.(*User) // Changed to *User
				if !ok {
					t.Fatalf("expected *User, got %T", model)
				}
				if savedUser.ID != userToSave.ID {
					t.Errorf("Save called with wrong user ID: expected %d, got %d", userToSave.ID, savedUser.ID)
				}
				return MockResult{RowsAffectedValue: 1}, nil
			},
		}
		ds := NewORMDataStore(mockRepo)

		err := ds.SaveUser(userToSave)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	// Test Case 2: Underlying Repository Error
	t.Run("RepoError", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("transaction failed")
		mockRepo := &MockRepository{
			SaveFn: func(model any) (sql.Result, error) {
				return nil, expectedErr
			},
		}
		ds := NewORMDataStore(mockRepo)

		err := ds.SaveUser(&User{ID: 11}) // Changed to *User

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error to wrap %v, got %v", expectedErr, err)
		}
	})

	// Test Case 3: Nil Repository
	t.Run("NilRepo", func(t *testing.T) {
		t.Parallel()
		ds := &ORMDataStore{Repo: nil}

		err := ds.SaveUser(&User{ID: 12})

		if err == nil {
			t.Fatal("expected an error for nil repository, got nil")
		}
		if err.Error() != "repository is nil" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestORMDataStore_DeleteUser(t *testing.T) {
	t.Parallel()

	// Test Case 1: Successful Deletion
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		deletedID := 20

		mockRepo := &MockRepository{
			DeleteFn: func(model any) (sql.Result, error) {
				user, ok := model.(*User) // Changed to *User
				if !ok {
					t.Fatalf("expected *User, got %T", model)
				}
				if user.ID != deletedID {
					t.Errorf("Delete called with wrong user ID: expected %d, got %d", deletedID, user.ID)
				}
				return MockResult{RowsAffectedValue: 1}, nil
			},
		}
		ds := NewORMDataStore(mockRepo)

		err := ds.DeleteUser(deletedID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	// Test Case 2: Underlying Repository Error
	t.Run("RepoError", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("permission denied")
		mockRepo := &MockRepository{
			DeleteFn: func(model any) (sql.Result, error) {
				return nil, expectedErr
			},
		}
		ds := NewORMDataStore(mockRepo)

		err := ds.DeleteUser(21)

		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error to wrap %v, got %v", expectedErr, err)
		}
	})

	// Test Case 3: Nil Repository
	t.Run("NilRepo", func(t *testing.T) {
		t.Parallel()
		ds := &ORMDataStore{Repo: nil}

		err := ds.DeleteUser(22)

		if err == nil {
			t.Fatal("expected an error for nil repository, got nil")
		}
		if err.Error() != "repository is nil" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

// --- Tests for APIDataStore ---

func TestAPIDataStore_GetUserByID(t *testing.T) {
	t.Parallel()
	ds := &APIDataStore{} // Changed to APIDataStore

	// Test Case 1: Found (ID 1)
	t.Run("Found", func(t *testing.T) {
		t.Parallel()
		user, err := ds.GetUserByID(1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.ID != 1 || user.Name != "Mock API User" {
			t.Errorf("unexpected user data: %+v", user)
		}
	})

	// Test Case 2: Not Found (ID != 1)
	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		user, err := ds.GetUserByID(2)

		if user != nil {
			t.Fatalf("expected nil user, got %+v", user)
		}
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		expectedErr := "API not implemented: user not found"
		if err.Error() != expectedErr {
			t.Errorf("unexpected error message: expected %q, got %q", expectedErr, err.Error())
		}
	})
}

func TestAPIDataStore_SaveUser(t *testing.T) {
	t.Parallel()
	ds := &APIDataStore{}

	err := ds.SaveUser(&User{ID: 1, Name: "Should Fail"})

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	expectedErr := "API not implemented: cannot save user"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: expected %q, got %q", expectedErr, err.Error())
	}
}

func TestAPIDataStore_DeleteUser(t *testing.T) {
	t.Parallel()
	ds := &APIDataStore{}

	err := ds.DeleteUser(1)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	expectedErr := "API not implemented: cannot delete user"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: expected %q, got %q", expectedErr, err.Error())
	}
}
