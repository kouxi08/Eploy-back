package repository

import (
	"context"
	"database/sql"

	"github.com/kouxi08/Eploy/internal/domain"
)

type ProjectRepository interface {
	GetProjectsByUserID(ctx context.Context, userId int) ([]domain.Project, error)
	CreateProject(ctx context.Context, tx *sql.Tx, project domain.Project, userId int) (int, error)
	CreateProjectEnvironments(ctx context.Context, tx *sql.Tx, projectId int, environments []domain.Environment) error
	CreateProjectWithEnvironments(ctx context.Context, project domain.Project, userId int) error
	GetProjectByID(ctx context.Context, id int, userId int) (domain.Project, error)
	DeleteProjectEnvironments(ctx context.Context, tx *sql.Tx, projectId int) error
	DeleteProject(ctx context.Context, tx *sql.Tx, projectId int, userId int) error
	DeleteProjectWithEnvironments(ctx context.Context, projectId int, userId int) error
}
