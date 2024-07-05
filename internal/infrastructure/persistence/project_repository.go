package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kouxi08/Eploy/internal/domain"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetProjectsByUserID(ctx context.Context, userId int) ([]domain.Project, error) {
	query := `
		SELECT 
			id,
			name,
			git_repo_url,
			domain,
			deployment_name,
		FROM 
			projects 
		WHERE
			user_id = ?`
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.GitRepoURL, &project.Domain, &project.DeploymentName); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		// プロジェクトが見つからない場合は空のスライスを返す
		return []domain.Project{}, nil
	}

	// 一覧画面なので環境変数は取得しない
	return projects, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, tx *sql.Tx, project domain.Project, userId int) (int, error) {
	query := `
		INSERT INTO 
			projects
			(
				name, 
				user_id, 
				git_repo_url, 
				domain,
				deployment_name, 
				dockerfile_dir, 
				port
			)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(
		ctx,
		query,
		project.Name,
		userId,
		project.GitRepoURL,
		project.Domain,
		project.DeploymentName,
		project.DockerfileDir,
		project.Port,
	)
	if err != nil {
		return 0, err
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(projectID), nil
}

func (r *ProjectRepository) CreateProjectEnvironments(ctx context.Context, tx *sql.Tx, projectID int, environments []domain.Environment) error {
	query := `
		INSERT INTO 
			project_env_vars
			(
				project_id, 
				env_key, 
				env_value
			)
		VALUES
			(?, ?, ?)`

	for _, env := range environments {
		_, err := tx.ExecContext(
			ctx,
			query,
			projectID,
			env.EnvKey,
			env.EnvValue,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectRepository) CreateProjectWithEnvironments(ctx context.Context, project domain.Project, userId int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	projectID, err := r.CreateProject(ctx, tx, project, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.CreateProjectEnvironments(ctx, tx, projectID, project.Environments)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id int, userId int) (domain.Project, error) {
	var project domain.Project
	query := `
		SELECT
			id,
			name,
			domain,
			git_repo_url,
			created_at
		FROM
			projects
		WHERE
			id = ? AND user_id = ?`

	err := r.db.QueryRowContext(ctx, query, id, userId).Scan(&project.ID, &project.Name, &project.Domain, &project.GitRepoURL, &project.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Project{}, fmt.Errorf("project not found")
		}
		return domain.Project{}, err
	}

	// 環境変数を取得する
	envQuery := `
		SELECT
			env_key,
			env_value
		FROM
			project_env_vars
		WHERE
			project_id = ?`
	rows, err := r.db.QueryContext(ctx, envQuery, project.ID)
	if err != nil {
		return domain.Project{}, err
	}
	defer rows.Close()

	var environments []domain.Environment
	for rows.Next() {
		var env domain.Environment
		if err := rows.Scan(&env.EnvKey, &env.EnvValue); err != nil {
			return domain.Project{}, err
		}
		environments = append(environments, env)
	}

	project.Environments = environments

	return project, nil
}

func (r *ProjectRepository) DeleteProjectEnvironments(ctx context.Context, tx *sql.Tx, projectId int) error {
	query := `
		DELETE FROM
			project_env_vars
		WHERE
			project_id = ?`
	_, err := tx.ExecContext(ctx, query, projectId)
	return err
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, tx *sql.Tx, projectId int, userId int) error {
	query := `
		DELETE FROM
			projects
		WHERE
			id = ? AND user_id = ?`
	_, err := tx.ExecContext(ctx, query, projectId, userId)
	return err
}

func (r *ProjectRepository) DeleteProjectWithEnvironments(ctx context.Context, projectId int, userId int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// プロジェクトの環境変数を削除
	err = r.DeleteProjectEnvironments(ctx, tx, projectId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// プロジェクトを削除
	err = r.DeleteProject(ctx, tx, projectId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
