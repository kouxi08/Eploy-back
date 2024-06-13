package main

import (
	"github.com/kouxi08/Eploy/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//サーバ起動
	server()
	// handler.LogsTest()
}

func server() {
	//インスタンス作成
	e := echo.New()

	//ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//podのログを取得(クエリパラメータ,podName="ポッド名")
	e.GET("/getpodlog", handler.GetPodLogHandler)

	//リソース追加処理へ
	e.POST("/", handler.CreateHandler)

	//kanikoのjobを起動する処理
	e.POST("/kaniko", handler.CreateKanikoHandler)

	//リソース削除処理へ
	e.PATCH("/", handler.DeleteHandler)

	e.GET("/", handler.GetMysqlPodLogHandler)
	e.GET("/dashboard", handler.GetDashboard)
	e.POST("createapp", handler.CreateAPP)

	e.Logger.Fatal(e.Start(":8088"))
}
