package orm

import (
	"fmt"
)

// Config holds the configuration options for the Liteforge database connection.
type Config struct {
	DriverName        string // The name of the database driver (e.g., "sqlite3", "postgres").
	DataSourceName    string // The connection string for the database.
	UseWriteAheadLogs bool   // Whether to enable Write Ahead Logs for sqlite
	EncryptAtRest     bool   // Whether to enable encryption at rest (SQLCipher for SQLite).
	EncryptionKey     string // The encryption key (if EncryptAtRest is true).  SHOULD NOT BE HARDCODED.
}

// OpenDB establishes a database connection based on the provided configuration and returns a Datastore.
func OpenDB(cfg Config) (*Datastore, error) {
	var adapter DBAdapter
	switch cfg.DriverName {
	case "sqlite3":
		adapter = &SQLiteAdapter{}
	case "postgres":
		adapter = &PostgresAdapter{}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DriverName)
	}

	db, err := adapter.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &Datastore{
		DB:      db,
		Adapter: adapter,
	}, nil
}
