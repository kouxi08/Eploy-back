package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kouxi08/Eploy/internal/domain"
	"github.com/kouxi08/Eploy/internal/interfaces/repository"
	"github.com/kouxi08/Eploy/pkg"
	"github.com/kouxi08/Eploy/pkg/cloudflare"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/kouxi08/Eploy/utils"
)

type ProjectUsecase struct {
	ProjectRepo   repository.ProjectRepository
	CloudflareApp *cloudflare.CloudflareAPP
	ConfigData    *utils.Config
	KubernetesApp *kubernetes.KubernetesApp
}

func NewProjectUsecase(repo repository.ProjectRepository, app *cloudflare.CloudflareAPP, config *utils.Config, kube *kubernetes.KubernetesApp) *ProjectUsecase {
	return &ProjectUsecase{
		ProjectRepo:   repo,
		CloudflareApp: app,
		ConfigData:    config,
		KubernetesApp: kube,
	}
}

func (u *ProjectUsecase) GetProjects(ctx context.Context, userId int) ([]domain.Project, error) {
	projects, err := u.ProjectRepo.GetProjectsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	for i, project := range projects {
		status, err := pkg.GetStatusResources(u.KubernetesApp, project.DeploymentName, "kaniko")
		if err != nil {
			return nil, err
		}
		projects[i].Status = status
	}
	return projects, nil
}

func (u *ProjectUsecase) CreateProject(ctx context.Context, project domain.Project, userId int) error {
	err := u.ProjectRepo.GetProjectName(ctx, userId, project.Name)
	if err != nil {
		return err
	}

	// プロジェクトの環境変数をEnvVarに変換する
	var envVars []kubernetes.EnvVar
	for _, env := range project.Environments {
		envVar := kubernetes.EnvVar{
			Name:  env.EnvKey,
			Value: env.EnvValue,
		}
		envVars = append(envVars, envVar)
	}

	Result, err := pkg.CreateResourceNames(&u.ConfigData.KubeManifest, project.Name, strconv.Itoa(project.Port))
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		// CreateResoucesを呼び出す
		err = pkg.CreateResouces(u.KubernetesApp, project.GitRepoURL, project.Name, Result, envVars)
		if err != nil {
			log.Printf("failed to create resources: %v", err)
			errChan <- err // エラーをチャネルに送信
		}
	}()

	// kanikoResultのDeploymentNameをプロジェクトのDeploymentNameに設定する
	project.Domain = Result.HostName
	project.DeploymentName = Result.DeploymentName

	// プロジェクトをリポジトリに保存する
	if err := u.ProjectRepo.CreateProjectWithEnvironments(ctx, project, userId); err != nil {
		// 何らかの理由で保存に失敗した場合は、削除処理を行う
		// deleteErr := pkg.DeleteResources(Result.DeploymentName)
		// if deleteErr != nil {
		// 	// 削除も失敗した場合はログなどで通知するなどの対応が必要です
		// 	return fmt.Errorf("failed to create project and failed to clean up resources: %v, delete error: %v", err, deleteErr)
		// }
		return fmt.Errorf("failed to create project: %v", err)
	}

	err = u.CloudflareApp.AddRecord(&u.ConfigData.DNSRecords, project.Name)
	if err != nil {
		return err
	}

	if err = <-errChan; err != nil {
		// エラーが発生していた場合は、削除処理を行う
		deleteErr := pkg.DeleteResources(u.KubernetesApp, &u.ConfigData.KubeManifest, Result.DeploymentName)
		if deleteErr != nil {

			return fmt.Errorf("failed to create project and failed to clean up resources: %v, delete error: %v", err, deleteErr)
		}
		return fmt.Errorf("failed to create resources: %v", err)
	}

	return nil
}

func (u *ProjectUsecase) GetProjectByID(ctx context.Context, id int, userId int) (domain.Project, error) {
	project, err := u.ProjectRepo.GetProjectByID(ctx, id, userId)
	if err != nil {
		return domain.Project{}, err
	}
	status, err := pkg.GetStatusResources(u.KubernetesApp, project.DeploymentName, "kaniko")
	if err != nil {
		return domain.Project{}, err
	}
	project.Status = status
	return project, nil
}

func (u *ProjectUsecase) GetProjectStatusByDeploymentName(ctx context.Context, deploymentName string) (string, error) {
	return pkg.GetStatusResources(u.KubernetesApp, deploymentName, "kaniko")
}

func (u *ProjectUsecase) DeleteProject(ctx context.Context, id int, userId int) error {
	// プロジェクト ID からプロジェクト情報を取得
	project, err := u.ProjectRepo.GetProjectByID(ctx, id, userId)
	if err != nil {
		return err
	}
	// DeleteResourcesを呼び出す
	err = pkg.DeleteResources(u.KubernetesApp, &u.ConfigData.KubeManifest, project.Name)
	if err != nil {
		return err
	}
	err = u.ProjectRepo.DeleteProjectWithEnvironments(ctx, id, userId)
	if err != nil {
		return err
	}

	err = u.CloudflareApp.DeleteRecord(&u.ConfigData.DNSRecords, project.Name)
	if err != nil {
		return err
	}

	return nil
}
