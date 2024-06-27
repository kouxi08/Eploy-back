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

	e.GET("/", handler.GetMysqlPodLogHandler)

	//リソース削除処理へ
	e.PATCH("/", handler.DeleteHandler)

	//リソース追加処理へ
	e.POST("/", handler.CreateHandler)

	//ダッシュボード一覧取得
	e.GET("/dashboard", handler.GetDashboard)

	//podのログを取得(クエリパラメータ,podName="ポッド名")
	e.GET("/getpodlog", handler.GetPodLogHandler)

	e.Logger.Fatal(e.Start(":8088"))
}
