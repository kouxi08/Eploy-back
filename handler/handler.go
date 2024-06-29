package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kouxi08/Eploy/pkg"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string `json:"message"`
}

// アプリケーションの作成
func CreateHandler(c echo.Context) error {
	requestData := new(kubernetes.RequestData)
	fmt.Println("Received JSON:", requestData)
	if err := c.Bind(requestData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	response := Response{
		Message: "Resources Create successfully",
	}

	err := pkg.CreateKanikoResouces(requestData.URL, requestData.Name, requestData.Port, requestData.EnvVars)
	if err != nil {
		response = Response{
			Message: "Resources Create failed",
		}
		log.Print(err)
	}
	return c.JSON(http.StatusOK, response)
}

// アプリケーションの削除
func DeleteHandler(c echo.Context) error {
	siteName := c.FormValue("name")
	response := Response{
		Message: "Resources Delete successfully",
	}
	err := pkg.DeleteResources(siteName)
	if err != nil {
		response = Response{
			Message: "Resources Delete failed",
		}
		log.Print(err)
	}
	return c.JSON(http.StatusOK, response)
}

// アプリケーションのログを取得
func GetPodLogHandler(c echo.Context) error {
	podName := c.QueryParam("podName")
	if podName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "The 'podname' query parameter is missing."})
	}
	log.Println(podName)
	// a,err := pkg.GetPodLog(podName)
	resultMessage, err := pkg.GetLogPodResources(podName)
	if err != nil {
		errStr := err.Error()
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errStr})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": resultMessage})
}

func GetMysqlPodLogHandler(c echo.Context) error {
	//  databaseの接続処理
	db, err := pkg.InitMysql()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// アクセス処理
	result, err := pkg.GetAccessLogs(db, "nginx.fast")
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, result)
}

func GetDashboard(c echo.Context) error {
	// userid := 1
	userid, err := strconv.Atoi(c.QueryParam("userid"))

	db, err := pkg.InitMysql()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	result, err := pkg.GetApp(db, userid)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}
