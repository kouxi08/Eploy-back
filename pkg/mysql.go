package pkg

import (
	  _ "github.com/go-sql-driver/mysql"
	  "database/sql"
)

func InitMysql() (db *sql.DB,err error) {
  db, err = sql.Open("mysql", "root:@(192.168.130.101:30010)/fluentd")
  if err != nil {
      return nil, err
  }
  return 
}

func GetLogs(db *sql.DB,url string) ([]LogsJSON,error){
  stmt, err := db.Prepare("SELECT * FROM nginx_access_logs where vhost = ?")
  if err != nil {
    panic(err.Error())
  }
  rows, err := stmt.Query(url)

  log := Logs{}
  var result []LogsJSON

  for rows.Next() {
    err := rows.Scan(
      &log.ID, &log.Time, &log.RemoteAddr, &log.XForwardedFor, &log.RequestID,
      &log.RemoteUser, &log.BytesSent, &log.RequestTime, &log.Status, &log.Vhost,
      &log.RequestProto, &log.Path, &log.RequestQuery, &log.RequestLength, &log.Duration,
      &log.Method, &log.HTTPReferrer, &log.HTTPUserAgent,
    )
    if err != nil {
      return nil,err
    }else {
      logJSON := convertToJSON(log)
      result = append(result,logJSON)
    }
  }
  return result ,nil
}