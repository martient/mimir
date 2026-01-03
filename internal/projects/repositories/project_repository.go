package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/taultek/mimir/internal/database"
	projModels "github.com/taultek/mimir/internal/projects/models"
)

// ProjectRepository handles project data operations
type ProjectRepository struct {
	db *database.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *database.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create inserts a new project
func (r *ProjectRepository) Create(project *projModels.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	project.UpdatedAt = time.Now()

	query := `
		INSERT INTO projects (id, name, path, opencode_port, project_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Conn.Exec(query,
		project.ID,
		project.Name,
		project.Path,
		project.opencodePort,
		project.ProjectType,
		project.CreatedAt.Unix(),
		project.UpdatedAt.Unix(),
	)
	return err
}

// Get retrieves a project by ID
func (r *ProjectRepository) Get(id string) (*projModels.Project, error) {
	var project projModels.Project

	query := `
		SELECT id, name, path, opencode_port, project_type, created_at, updated_at
		FROM projects
		WHERE id = ?
	`
	err := r.db.Conn.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Path,
		&project.opencodePort,
		&project.ProjectType,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// GetByName retrieves a project by name
func (r *ProjectRepository) GetByName(name string) (*projModels.Project, error) {
	var project projModels.Project

	query := `
		SELECT id, name, path, opencode_port, project_type, created_at, updated_at
		FROM projects
		WHERE name = ?
	`
	err := r.db.Conn.QueryRow(query, name).Scan(
		&project.ID,
		&project.Name,
		&project.Path,
		&project.opencodePort,
		&project.ProjectType,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// List retrieves all projects
func (r *ProjectRepository) List() ([]projModels.Project, error) {
	query := `
		SELECT id, name, path, opencode_port, project_type, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`
	rows, err := r.db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []projModels.Project
	for rows.Next() {
		var project projModels.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Path,
			&project.opencodePort,
			&project.ProjectType,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

// Update updates a project
func (r *ProjectRepository) Update(project *projModels.Project) error {
	query := `
		UPDATE projects
		SET name = ?, path = ?, opencode_port = ?, project_type = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Conn.Exec(query,
		project.Name,
		project.Path,
		project.opencodePort,
		project.ProjectType,
		time.Now().Unix(),
		project.ID,
	)
	return err
}

// Delete removes a project
func (r *ProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := r.db.Conn.Exec(query, id)
	return err
}
