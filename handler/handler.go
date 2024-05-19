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
