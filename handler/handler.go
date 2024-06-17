package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kouxi08/Eploy/pkg"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/labstack/echo/v4"
)

// アプリケーションの作成
func CreateHandler(c echo.Context) error {

	siteName := c.FormValue("name")
	targetPort := c.FormValue("port")

	pkg.CreateResources(siteName, targetPort)

	return c.String(http.StatusOK, "Resources added successfully")
}

// アプリケーションの削除
func DeleteHandler(c echo.Context) error {

	siteName := c.FormValue("name")

	pkg.DeleteResources(siteName)

	return c.String(http.StatusOK, "Resources  delete successfully")
}

// Kanikoの処理を作成
func CreateKanikoHandler(c echo.Context) error {
	//envファイルを受け渡すために構造体を引っ張ってきてる(他にいい方法があるはず)
	requestData := new(kubernetes.RequestData)
	println("Received JSON:", requestData)
	if err := c.Bind(requestData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	pkg.CreateKanikoResouces(requestData.URL, requestData.Name, requestData.Port, requestData.EnvVars)

	return c.String(http.StatusOK, "Job create successfully")
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

func GetDashboard(c echo.Context)error {
	userid := 1
	db, err := pkg.InitMysql()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	result,err := pkg.GetApp(db, userid)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}

// アプリを登録する時の処理、jsonで以下形式で登録
// {
//     "appName": "nginx",
//     "domain": "ru-ru.kouxi.com",
//     "gitURL": "https://github.com/kouxi08/pixivbot",
//     "deploymentName": "nginx-deployment"
// }
func CreateApp(c echo.Context)error{
	userID := 1 // 仮にuserIDは静的に設定

	data, err := pkg.BindData(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	db, err := pkg.InitMysql()
	if err != nil {
		log.Println("Database initialization failed:", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = pkg.InsertApp(db, data.AppName, userID, data.Domain, data.GitURL, data.DeploymentName)
	if err != nil {
		log.Println("Error inserting app data:", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "success"})
}
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test (c echo.Context)error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}