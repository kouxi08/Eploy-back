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

	return c.String(http.StatusOK, "Resources added successfully")
}

func CreateKanikoHandler(c echo.Context) error {
	pkg.CreateKanikoResouces()
	return c.String(http.StatusOK, "")
}

func DeleteHandler(c echo.Context) error {
	config, _ := config.LoadConfig("config.json")

	siteName := c.FormValue("name")

	deploymentName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.DeploymentName)
	serviceName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.ServiceName)
	ingressName := fmt.Sprintf("%s%s", siteName, config.KubeConfig.IngressName)

	pkg.DeleteResources(deploymentName, serviceName, ingressName)

	return c.String(http.StatusOK, "Resources  delete successfully")
}

func GetPodLogHandler(c echo.Context) error {
	podName := c.QueryParam("podName")
	if podName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "The 'podname' query parameter is missing."})
	}
	fmt.Println(podName)
	// a,err := pkg.GetPodLog(podName)
	resultMessage, err := pkg.GetLogPodResources(podName)
	if err != nil {
		errStr := err.Error()
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errStr})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": resultMessage})
}
