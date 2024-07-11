package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kouxi08/Eploy/config"
	projectRepo "github.com/kouxi08/Eploy/internal/infrastructure/persistence"
	projectHandler "github.com/kouxi08/Eploy/internal/interfaces/handler"
	customMiddleware "github.com/kouxi08/Eploy/internal/middleware"
	projectUsecase "github.com/kouxi08/Eploy/internal/usecase"
	"github.com/kouxi08/Eploy/pkg/firebase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//インスタンス作成
	e := echo.New()
	config.Env()

	firebaseApp, err := firebase.InitFirebaseApp()
	if err != nil {
		log.Fatalf("failed to initialize Firebase app: %v", err)
		os.Exit(1)
	}

	message := os.Getenv("MYSQL_URL")
	db, err := sql.Open("mysql", message)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepository := projectRepo.NewUserRepository(db)

	//ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(customMiddleware.AuthMiddleware(firebaseApp, userRepository))

	projectRepository := projectRepo.NewProjectRepository(db)
	projectUsecase := projectUsecase.NewProjectUsecase(projectRepository)
	projectHandler := projectHandler.NewProjectHandler(projectUsecase)

	e.GET("/projects", projectHandler.GetProjects)
	e.POST("/projects", projectHandler.CreateProject)
	e.GET("/projects/:id", projectHandler.GetProjectByID)
	e.GET("/projects/:deployment_name/status", projectHandler.GetProjectStatusByDeploymentName)
	e.GET("/projects/:deployment_name/logs", projectHandler.GetProjectLogsByDeploymentName)
	e.DELETE("/projects/:id", projectHandler.DeleteProject)

	e.Logger.Fatal(e.Start(":8088"))
}
