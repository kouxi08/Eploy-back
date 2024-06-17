package pkg
import (
	"github.com/labstack/echo/v4"
	"log"
)

// AppData はリクエストデータの構造体です。
type AppData struct {
	AppName        string `json:"appName"`
	Domain         string `json:"domain"`
	GitURL         string `json:"gitURL"`
	DeploymentName string `json:"deploymentName"`
}

// BindData はEchoのコンテキストからAppDataをバインドします。
func BindData(c echo.Context) (*AppData, error) {
	var data AppData
	if err := c.Bind(&data); err != nil {
		log.Println("Error binding data:", err)
		return nil, err
	}
	return &data, nil
}