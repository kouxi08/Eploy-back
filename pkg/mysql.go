package pkg

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// mysqlの初期の接続処理
func InitMysql() (db *sql.DB, err error) {
	// .envからmysqlのurlを取得
	err = godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
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
	// jsonに変換
	result, err := ConvertToJSON(rows)
	if err != nil {
		return nil, err
	}
	return result, nil
}
