package pkg

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kouxi08/Eploy/utils"
)
// type App struct {
// 	ID              int    `json:"id"`
// 	ApplicationName string `json:"application_name"`
// 	Domain          string `json:"domain"`
// 	GithubURL       string `json:"github_url"`
// 	Status          string `json:"status"`
// }
// // Response structure
// type Response struct {
// 	Sites []App `json:"sites"`
// }

// mysqlの初期の接続処理
func InitMysql() (db *sql.DB, err error) {
	// .envからmysqlのurlを取得
	utils.Env()
	message := os.Getenv("MYSQL_URL")
	// mysql接続
	db, err = sql.Open("mysql", message)
	if err != nil {
		return nil, err
	}
	return
}

// mysqlからlogを取得
func GetAccessLogs(db *sql.DB, url string) ([]LogsJSON, error) {
	// アクセスurlを絞る
	stmt, err := db.Prepare("SELECT * FROM nginx_access_logs where vhost = ?")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(url)
	if err != nil {
		return nil, err
	}
	// jsonに変換
	result, err := ConvertToJSON(rows)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetApp(db *sql.DB,userid int) (*Response,error)  {
	// app側からuseridを取得してくる　注データベースの設計前に作成しているため変更がいるかも
	stmt, err := db.Prepare("SELECT * FROM app where user_id = ?")
	if err != nil {
		return nil,err
	}
	// 
	rows, err := stmt.Query(userid)
	result ,err:= ConvertToJSONDs(rows)
	if err != nil {
		return nil,err
	}
	return  result,nil
}
