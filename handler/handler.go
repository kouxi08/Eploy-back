package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kouxi08/Eploy/config"
	"github.com/kouxi08/Eploy/pkg"
	"github.com/kouxi08/Eploy/pkg/kubernetes"
	"github.com/labstack/echo/v4"
)

// アプリケーションの作成
func CreateHandler(c echo.Context) error {
	config, _ := config.LoadConfig("config.json")

	siteName := c.FormValue("name")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)
	hostName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.HostName)

	pkg.CreateResources(siteName, deploymentName, serviceName, ingressName, hostName)

	return c.String(http.StatusOK, "Resources added successfully")
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

	pkg.CreateKanikoResouces(requestData.URL, requestData.Name, requestData.EnvVars)
	return c.String(http.StatusOK, "Job create successfully")
}

// アプリケーションの削除
func DeleteHandler(c echo.Context) error {
	config, _ := config.LoadConfig("config.json")

	siteName := c.FormValue("name")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)

	pkg.DeleteResources(deploymentName, serviceName, ingressName)

	return c.String(http.StatusOK, "Resources  delete successfully")
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

func GetMysqlPodLogHandler(c echo.Context) error{
	//  databaseの接続処理
	db,err := pkg.InitMysql()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// アクセス処理
	result,err := pkg.GetAccessLogs(db,"nginx.fast")
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, result)
}