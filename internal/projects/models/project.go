package projects

import (
	"time"
)

// Project represents a registered project
type Project struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Path         string    `db:"path"`
	opencodePort int       `db:"opencode_port"`
	ProjectType  string    `db:"project_type"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// ProjectType constants
const (
	TypeGo         = "golang"
	TypePython     = "python"
	TypeTypeScript = "typescript"
	TypeRust       = "rust"
	TypeUnknown    = "unknown"
)
