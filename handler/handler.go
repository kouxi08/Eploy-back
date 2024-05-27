package handler

import (
	"fmt"
	"net/http"

	"github.com/kouxi08/Eploy/config"
	"github.com/kouxi08/Eploy/pkg"
	"github.com/labstack/echo/v4"
)

func CreateHandler(c echo.Context) error {
	config, _ := config.LoadConfig("config.json")

	siteName := c.FormValue("name")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)
	hostName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.HostName)

	pkg.CreateResources(siteName, deploymentName, serviceName, ingressName, hostName)

	return c.String(http.StatusOK, "Record added successfully")
}

func DeleteHandler(c echo.Context) error {
	config, _ := config.LoadConfig("config.json")

	siteName := c.FormValue("name")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)

	pkg.DeleteResources(deploymentName, serviceName, ingressName)

	return c.String(http.StatusOK, "Record delete successfully")
}

func GetPodLogHandler(c echo.Context)error {
	podName := c.QueryParam("podName")
	if(podName == "") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "The 'podname' query parameter is missing."})
	}
	fmt.Println(podName)
	// a,err := pkg.GetPodLog(podName) 
	a,err := pkg.GetLogPodResources(podName) 
	if err != nil {
		errStr := err.Error()
		return c.JSON(http.StatusBadRequest, map[string]string{"error":errStr})
	}
	return c.JSON(http.StatusOK, map[string]string{"message":a})
}

func LogsTest(c echo.Context) error{
	db,err := pkg.InitMysql()
	if err != nil {
		fmt.Println("err init database ")
	}
	fmt.Print(db)
	result,err := pkg.GetLogs(db,"nginx.fast")
	if err != nil {
		fmt.Print("accesslog.GetLogs:")
		fmt.Println(err)	}
	fmt.Print(result)
	return c.JSON(http.StatusOK, result)
}