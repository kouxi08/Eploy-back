package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kouxi08/Eploy/internal/domain"
	"github.com/kouxi08/Eploy/internal/interfaces/repository"
	"github.com/kouxi08/Eploy/pkg"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
)

type ProjectUsecase struct {
	ProjectRepo repository.ProjectRepository
}

func NewProjectUsecase(repo repository.ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{
		ProjectRepo: repo,
	}
}

func (u *ProjectUsecase) GetProjects(ctx context.Context, userId int) ([]domain.Project, error) {
	projects, err := u.ProjectRepo.GetProjectsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	for i, project := range projects {
		status, err := pkg.GetStatusResources(project.DeploymentName)
		if err != nil {
			return nil, err
		}
		projects[i].Status = status
	}
	return projects, nil
}
func (u *ProjectUsecase) CreateProject(ctx context.Context, project domain.Project, userId int) error {
	// プロジェクトの環境変数をEnvVarに変換する
	var envVars []kubernetes.EnvVar
	for _, env := range project.Environments {
		envVar := kubernetes.EnvVar{
			Name:  env.EnvKey,
			Value: env.EnvValue,
		}
		envVars = append(envVars, envVar)
	}

	// CreateKanikoResoucesを呼び出す
	kanikoResult, err := pkg.CreateKanikoResouces(project.GitRepoURL, project.Name, strconv.Itoa(project.Port), envVars)
	if err != nil {
		return err
	}

	// kanikoResultのDeploymentNameをプロジェクトのDeploymentNameに設定する
	project.Domain = kanikoResult.HostName
	project.DeploymentName = kanikoResult.DeploymentName

	// プロジェクトをリポジトリに保存する
	if err := u.ProjectRepo.CreateProjectWithEnvironments(ctx, project, userId); err != nil {
		// 何らかの理由で保存に失敗した場合は、削除処理を行う
		deleteErr := pkg.DeleteResources(kanikoResult.DeploymentName)
		if deleteErr != nil {
			// 削除も失敗した場合はログなどで通知するなどの対応が必要です
			return fmt.Errorf("failed to create project and failed to clean up resources: %v, delete error: %v", err, deleteErr)
		}
		return fmt.Errorf("failed to create project: %v", err)
	}

	return nil
}

func (u *ProjectUsecase) GetProjectByID(ctx context.Context, id int, userId int) (domain.Project, error) {
	project, err := u.ProjectRepo.GetProjectByID(ctx, id, userId)
	if err != nil {
		return domain.Project{}, err
	}
	status, err := pkg.GetStatusResources(project.DeploymentName)
	if err != nil {
		return domain.Project{}, err
	}
	project.Status = status
	return project, nil
}

func (u *ProjectUsecase) GetProjectStatusByDeploymentName(ctx context.Context, deploymentName string) (string, error) {
	return pkg.GetStatusResources(deploymentName)
}

func (u *ProjectUsecase) DeleteProject(ctx context.Context, id int, userId int) error {
	// プロジェクト ID からプロジェクト情報を取得
	project, err := u.ProjectRepo.GetProjectByID(ctx, id, userId)
	if err != nil {
		return err
	}
	// DeleteResourcesを呼び出す
	err = pkg.DeleteResources(project.Name)
	if err != nil {
		return err
	}
	err = u.ProjectRepo.DeleteProjectWithEnvironments(ctx, id, userId)
	if err != nil {
		return err
	}
	return nil
}
