package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mutecomm/go-sqlcipher"

	"github.com/taultek/mimir/internal/config"
)

// DB wraps sql.DB with encryption
type DB struct {
	Conn *sql.DB
}

// Init initializes database with encryption
func Init(cfg config.DatabaseConfig) (*DB, error) {
	// Expand ~ in path
	path := expandPath(cfg.Path)

	// Ensure database directory exists
	if len(path) > 0 && path[0] == '/' {
		dir := path[:len(path)-1]
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	} else if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = homeDir + path[1:]
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Build connection string with encryption key
	// sqlcipher requires key in connection string
	dsn := fmt.Sprintf("file:%s?_pragma_key=%s", path, cfg.EncryptionKey)

	// Open database
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better performance
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set busy timeout
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	return &DB{conn: db}, nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	// Create sessions table
	if _, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			agent_type TEXT NOT NULL,
			worktree_path TEXT,
			branch_name TEXT,
			status TEXT NOT NULL DEFAULT 'created',
			metadata TEXT,
			created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		);
	`); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	// Create projects table
	if _, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			path TEXT NOT NULL,
			opencode_port INTEGER NOT NULL,
			project_type TEXT,
			created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		);
	`); err != nil {
		return fmt.Errorf("failed to create projects table: %w", err)
	}

	// Create worktrees table
	if _, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS worktrees (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			path TEXT NOT NULL,
			branch_name TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
			deleted_at INTEGER,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		);
	`); err != nil {
		return fmt.Errorf("failed to create worktrees table: %w", err)
	}

	// Create indexes
	if _, err := db.conn.Exec("CREATE INDEX IF NOT EXISTS idx_sessions_project ON sessions(project_id)"); err != nil {
		return fmt.Errorf("failed to create sessions index: %w", err)
	}

	if _, err := db.conn.Exec("CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status)"); err != nil {
		return fmt.Errorf("failed to create sessions status index: %w", err)
	}

	if _, err := db.conn.Exec("CREATE INDEX IF NOT EXISTS idx_worktrees_session ON worktrees(session_id)"); err != nil {
		return fmt.Errorf("failed to create worktrees session index: %w", err)
	}

	return nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return homeDir + path[1:]
	}
	return path
}
