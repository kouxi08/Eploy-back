package main

import (
	"fmt"
	"net/http"

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

	e.GET("/dashboard", dashboardHandler)

	// e.GET("/dashboard", handler.GetDashboard)
	e.POST("/createapp", handler.CreateApp)

	e.Logger.Fatal(e.Start(":8088"))
}

func dashboardHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "Bearer 1" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}
	fmt.Print("aaaa")
	// ダッシュボードデータを返す
	data := map[string]interface{}{
		"status": "success",
		"data":   "Here is your dashboard data.",
	}

	return c.JSON(http.StatusOK, data)
}
