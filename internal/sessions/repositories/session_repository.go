package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/taultek/mimir/internal/database"
	sessModels "github.com/taultek/mimir/internal/sessions/models"
)

// SessionRepository handles session data operations
type SessionRepository struct {
	db *database.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *database.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create inserts a new session
func (r *SessionRepository) Create(session *sessModels.Session) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	session.UpdatedAt = time.Now()

	query := `
		INSERT INTO sessions (id, project_id, agent_type, worktree_path, branch_name, status, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Conn.Exec(query,
		session.ID,
		session.ProjectID,
		session.AgentType,
		session.WorktreePath,
		session.BranchName,
		session.Status,
		session.Metadata,
		session.CreatedAt.Unix(),
		session.UpdatedAt.Unix(),
	)
	return err
}

// Get retrieves a session by ID
func (r *SessionRepository) Get(id string) (*sessModels.Session, error) {
	var session sessModels.Session

	query := `
		SELECT id, project_id, agent_type, worktree_path, branch_name, status, metadata, created_at, updated_at
		FROM sessions
		WHERE id = ?
	`
	err := r.db.Conn.QueryRow(query, id).Scan(
		&session.ID,
		&session.ProjectID,
		&session.AgentType,
		&session.WorktreePath,
		&session.BranchName,
		&session.Status,
		&session.Metadata,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// List retrieves all sessions
func (r *SessionRepository) List() ([]sessModels.Session, error) {
	query := `
		SELECT id, project_id, agent_type, worktree_path, branch_name, status, metadata, created_at, updated_at
		FROM sessions
		ORDER BY created_at DESC
	`
	rows, err := r.db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []sessModels.Session
	for rows.Next() {
		var session sessModels.Session
		err := rows.Scan(
			&session.ID,
			&session.ProjectID,
			&session.AgentType,
			&session.WorktreePath,
			&session.BranchName,
			&session.Status,
			&session.Metadata,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// UpdateStatus updates a session's status
func (r *SessionRepository) UpdateStatus(id string, status string) error {
	query := `
		UPDATE sessions
		SET status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Conn.Exec(query, status, time.Now().Unix(), id)
	return err
}

// Delete removes a session
func (r *SessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Conn.Exec(query, id)
	return err
}

// ListByProject retrieves sessions for a specific project
func (r *SessionRepository) ListByProject(projectID string) ([]sessModels.Session, error) {
	query := `
		SELECT id, project_id, agent_type, worktree_path, branch_name, status, metadata, created_at, updated_at
		FROM sessions
		WHERE project_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Conn.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions for project: %w", err)
	}
	defer rows.Close()

	var sessionsList []sessModels.Session
	for rows.Next() {
		var session sessModels.Session
		err := rows.Scan(
			&session.ID,
			&session.ProjectID,
			&session.AgentType,
			&session.WorktreePath,
			&session.BranchName,
			&session.Status,
			&session.Metadata,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessionsList = append(sessionsList, session)
	}

	return sessionsList, rows.Err()
}
