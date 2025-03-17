package orm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // For SQLite
)

// Config holds the configuration options for the Liteforge database connection.
type Config struct {
	DriverName        string // The name of the database driver (e.g., "sqlite3", "postgres").
	DataSourceName    string // The connection string for the database.
	UseWriteAheadLogs bool   // Whether to enable Write Ahead Logs for sqlite
	EncryptAtRest     bool   // Whether to enable encryption at rest (SQLCipher for SQLite).
	EncryptionKey     string // The encryption key (if EncryptAtRest is true).  SHOULD NOT BE HARDCODED.
}

// OpenDB establishes a database connection based on the provided configuration.
func OpenDB(cfg Config) (*sql.DB, error) {
	dbDirPerm := 0755 //Directory permission
	dbPath := cfg.DataSourceName
	dbDir := filepath.Dir(dbPath)

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		//creating directory
		if err := os.MkdirAll(dbDir, os.FileMode(dbDirPerm)); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	} 

	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	//Set up wal mode, if needed.
	if cfg.UseWriteAheadLogs {
		_, err = db.Exec("PRAGMA journal_mode=WAL;")
		if err != nil {
			return nil, fmt.Errorf("Failed to switch to WAL mode: %w", err)
		}
	}
	// Test the connection

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
