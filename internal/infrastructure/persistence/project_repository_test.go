package persistence

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	testHelper "github.com/kouxi08/Eploy/test"

	"github.com/kouxi08/Eploy/internal/domain"
)

func TestCreateProject(t *testing.T) {
	repo := NewProjectRepository(db)
	ctx := context.Background()

	// テスト用のユーザー ID
	userId, err := testHelper.CreateUser(ctx, db)

	// テスト用のプロジェクトデータ
	project := domain.Project{
		Name:           "Test Project",
		GitRepoURL:     "https://github.com/example/test",
		Domain:         "example.com",
		DeploymentName: "test-deployment",
		DockerfileDir:  "/path/to/dockerfile",
		Port:           8080,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("トランザクション開始エラー: %v", err)
	}
	defer tx.Rollback()

	// プロジェクトの作成をテスト
	projectID, err := repo.CreateProject(ctx, tx, project, userId)
	if err != nil {
		t.Fatalf("プロジェクト作成のテストに失敗しました: %v", err)
	}

	// プロジェクト ID を取得して、取得できることを確認
	if projectID <= 0 {
		t.Fatalf("期待するプロジェクト ID が取得されませんでした。取得された ID: %d", projectID)
	}

	fmt.Printf("プロジェクト ID: %d\n", projectID)

	// プロジェクトが作成されたことを確認
	storedProjects, err := repo.GetProjectByID(ctx, projectID, userId)
	if err != nil {
		t.Fatalf("プロジェクトの取得に失敗しました: %v", err)
	}

	// storedProjectsとprojectが一致することを確認
	if storedProjects.ID != projectID {
		t.Fatalf("プロジェクト ID が一致しません。期待する ID: %d, 取得された ID: %d", projectID, storedProjects.ID)
	}

	if storedProjects.Name != project.Name {
		t.Fatalf("プロジェクト名が一致しません。期待する名前: %s, 取得された名前: %s", project.Name, storedProjects.Name)
	}

	if storedProjects.GitRepoURL != project.GitRepoURL {
		t.Fatalf("Git リポジトリ URL が一致しません。期待する URL: %s, 取得された URL: %s", project.GitRepoURL, storedProjects.GitRepoURL)
	}

	if storedProjects.Domain != project.Domain {
		t.Fatalf("ドメインが一致しません。期待するドメイン: %s, 取得されたドメイン: %s", project.Domain, storedProjects.Domain)
	}

	if storedProjects.DeploymentName != project.DeploymentName {
		t.Fatalf("デプロイメント名が一致しません。期待するデプロイメント名: %s, 取得されたデプロイメント名: %s", project.DeploymentName, storedProjects.DeploymentName)
	}

	// コミットしてトランザクションを終了
	err = tx.Commit()
	if err != nil {
		t.Fatalf("トランザクションのコミットに失敗しました: %v", err)
	}
}
