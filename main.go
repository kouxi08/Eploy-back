package main

import (
	"github.com/kouxi08/Eploy/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//インスタンス作成
	e := echo.New()

	//ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//podのログを取得(クエリパラメータ,podName="ポッド名")
	e.GET("/getpodlog", handler.GetPodLogHandler)

	e.GET("/", handler.GetMysqlPodLogHandler)

	//リソース削除処理へ
	e.PATCH("/", handler.DeleteHandler)

	//リソース追加処理へ
	e.POST("/", handler.CreateHandler)

	e.GET("/dashboard", handler.GetDashboard)

	e.POST("/createapp", handler.CreateApp)

	e.Logger.Fatal(e.Start(":8088"))
}
