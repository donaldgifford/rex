package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Database wraps the SQL database connection
type Database struct {
	conn *sql.DB
	path string
}

// Open opens or creates the SQLite database at the specified path
func Open(path string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	// Open SQLite connection
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Set pragmas for performance and data integrity
	pragmas := []string{
		"PRAGMA journal_mode=WAL",       // Write-Ahead Logging for better concurrency
		"PRAGMA synchronous=NORMAL",      // Balance between safety and speed
		"PRAGMA foreign_keys=ON",         // Enable foreign key constraints
		"PRAGMA busy_timeout=5000",       // Wait up to 5 seconds if database is locked
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			conn.Close()
			return nil, fmt.Errorf("set pragma: %w", err)
		}
	}

	db := &Database{
		conn: conn,
		path: path,
	}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// initSchema creates all tables and indexes if they don't exist
func (db *Database) initSchema() error {
	// Execute schema
	if _, err := db.conn.Exec(Schema); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}

	// Set schema version
	_, err := db.conn.Exec(
		"INSERT OR REPLACE INTO metadata (key, value, updated_at) VALUES (?, ?, datetime('now'))",
		"schema_version", "1",
	)
	if err != nil {
		return fmt.Errorf("set schema version: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// GetLastRebuild returns the timestamp of the last cache rebuild
func (db *Database) GetLastRebuild() (time.Time, error) {
	var timestamp string
	err := db.conn.QueryRow(
		"SELECT value FROM metadata WHERE key = ?",
		"last_rebuild",
	).Scan(&timestamp)

	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}

	if err != nil {
		return time.Time{}, err
	}

	// Parse RFC3339 timestamp
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse timestamp: %w", err)
	}

	return t, nil
}

// SetLastRebuild updates the last rebuild timestamp to now
func (db *Database) SetLastRebuild() error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.conn.Exec(
		"INSERT OR REPLACE INTO metadata (key, value, updated_at) VALUES (?, ?, datetime('now'))",
		"last_rebuild", now,
	)
	return err
}

// GetMetadata retrieves a metadata value by key
func (db *Database) GetMetadata(key string) (string, error) {
	var value string
	err := db.conn.QueryRow(
		"SELECT value FROM metadata WHERE key = ?",
		key,
	).Scan(&value)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return value, err
}

// SetMetadata sets a metadata key-value pair
func (db *Database) SetMetadata(key, value string) error {
	_, err := db.conn.Exec(
		"INSERT OR REPLACE INTO metadata (key, value, updated_at) VALUES (?, ?, datetime('now'))",
		key, value,
	)
	return err
}

// Conn returns the underlying *sql.DB connection
// This allows other packages to execute custom queries
func (db *Database) Conn() *sql.DB {
	return db.conn
}

// Path returns the database file path
func (db *Database) Path() string {
	return db.path
}
