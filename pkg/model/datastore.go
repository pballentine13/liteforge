package model

import (
	"database/sql"
	"errors"
	"fmt"
)

// DataStore defines the application-specific interface for data access.
// This abstracts the underlying data source (e.g., ORM, API, cache).
type DataStore interface {
	GetUserByID(id int) (*User, error)
	SaveUser(user *User) error
	DeleteUser(id int) error
}

// ORMDataStore is a concrete implementation of DataStore that uses the ORMRepository.
type ORMDataStore struct {
	Repo Repository
}

// NewORMDataStore creates a new ORMDataStore instance.
func NewORMDataStore(repo Repository) *ORMDataStore {
	return &ORMDataStore{Repo: repo}
}

// GetUserByID retrieves a User by their ID using the ORMRepository.
func (ds *ORMDataStore) GetUserByID(id int) (*User, error) {
	if ds.Repo == nil {
		return nil, errors.New("repository is nil")
	}

	user := &User{}
	if err := ds.Repo.FindByID(user, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil user, nil error if not found
		}
		return nil, fmt.Errorf("failed to find user by ID %d: %w", id, err)
	}
	return user, nil
}

// SaveUser saves a User (insert or update) using the ORMRepository.
func (ds *ORMDataStore) SaveUser(user *User) error {
	if ds.Repo == nil {
		return errors.New("repository is nil")
	}

	_, err := ds.Repo.Save(user)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// DeleteUser deletes a User by their ID.
// It creates a temporary User struct with the ID set for the deletion.
func (ds *ORMDataStore) DeleteUser(id int) error {
	if ds.Repo == nil {
		return errors.New("repository is nil")
	}

	userToDelete := &User{ID: id}
	_, err := ds.Repo.Delete(userToDelete)
	if err != nil {
		return fmt.Errorf("failed to delete user with ID %d: %w", id, err)
	}
	return nil
}

// APIDataStore is a mock implementation of DataStore for an external API.
type APIDataStore struct{}

// GetUserByID returns a dummy user or an error.
func (ds *APIDataStore) GetUserByID(id int) (*User, error) {
	if id == 1 {
		return &User{ID: 1, Name: "Mock API User", Age: 30}, nil
	}
	return nil, errors.New("API not implemented: user not found")
}

// SaveUser returns a specific error.
func (ds *APIDataStore) SaveUser(user *User) error {
	return errors.New("API not implemented: cannot save user")
}

// DeleteUser returns a specific error.
func (ds *APIDataStore) DeleteUser(id int) error {
	return errors.New("API not implemented: cannot delete user")
}
