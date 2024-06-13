package pkg

import (
	"fmt"
	"database/sql"
)

type App struct {
	ID              int    `json:"id"`
	ApplicationName string `json:"application_name"`
	Domain          string `json:"domain"`
	GithubURL       string `json:"github_url"`
	Status          string `json:"status"`
}
// Response structure
type Response struct {
	Sites []App `json:"sites"`
}

func ConvertToJSONDs(rows *sql.Rows) (*Response,error){
	// mysqlから取得してきたものをjson形式に治す
	var apps []App
	var dname string
	var uid string
	for rows.Next() {
		var app App
		if err := rows.Scan(&app.ID, &app.ApplicationName,&uid ,&app.Domain, &app.GithubURL, &dname); err != nil {
			return nil,err
		}
		status,err := GetStatusResources(dname)
		if err != nil{
			fmt.Print(status)
		}
        app.Status = status
		apps = append(apps, app)
	}
	response := Response{Sites: apps}
	return &response,nil
}