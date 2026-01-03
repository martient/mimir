package sessions

import (
	"time"
)

// Session represents an opencode agent session
type Session struct {
	ID           string    `db:"id"`
	ProjectID    string    `db:"project_id"`
	AgentType    string    `db:"agent_type"`
	WorktreePath string    `db:"worktree_path"`
	BranchName   string    `db:"branch_name"`
	Status       string    `db:"status"`
	Metadata     string    `db:"metadata"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// SessionStatus constants
const (
	StatusCreated   = "created"
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusError     = "error"
	StatusCancelled = "cancelled"
)
