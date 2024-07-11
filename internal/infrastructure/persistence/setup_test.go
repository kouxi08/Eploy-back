package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB
var pool *dockertest.Pool
var resource *dockertest.Resource

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")

	pwd, _ := os.Getwd()

	// MySQL コンテナの起動
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "latest",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=pass",
			"MYSQL_DATABASE=testdb",
		},
		Mounts: []string{
			// 絶対パスで指定しろと言われるので、絶対パスで指定
			pwd + "/../../../init.sql:/docker-entrypoint-initdb.d/init.sql",
		},
	}, func(config *docker.HostConfig) {
		// コンテナのポートを適切にマップする
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start MySQL container: %s", err)
	}

	// MySQL コンテナの接続待機
	hostPort := resource.GetPort("3306/tcp")
	dsn := fmt.Sprintf("root:pass@tcp(localhost:%s)/testdb?parseTime=true", hostPort)
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Failed to connect to MySQL container: %s", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()

}
