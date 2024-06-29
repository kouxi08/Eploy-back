package pkg

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kouxi08/Eploy/config"
)

// mysqlの初期の接続処理
func InitMysql() (db *sql.DB, err error) {
	// .envからmysqlのurlを取得
	config.Env()
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

func GetApp(db *sql.DB, userid int) (*Response, error) {
	// app側からuseridを取得してくる　注データベースの設計前に作成しているため変更がいるかも
	stmt, err := db.Prepare("SELECT * FROM app where user_id = ?")
	if err != nil {
		return nil, err
	}
	//
	rows, err := stmt.Query(userid)
	if err != nil {
		return nil, err
	}
	result, err := ConvertToJSONDs(rows)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertApp(db *sql.DB, appName string, userid int, domain string, gitURL string, deploymentName string) error {
	stmt, err := db.Prepare("INSERT INTO app(application_name,user_id,domain,github_url,deployment_name) VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(appName, userid, domain, gitURL, deploymentName)
	if err != nil {
		return err
	}
	// 成功時
	return nil
}

// appを削除する
func DeleteApp(db *sql.DB, deploymentName string) error {
	stmt, err := db.Prepare("DELETE FROM app WHERE deployment_name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(deploymentName)
	if err != nil {
		return err
	}
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffect == 0 {
		log.Println("no rows deleted")
		return nil
	}
	// 成功時
	return nil
}
