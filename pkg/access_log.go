package pkg

import (
	"database/sql"
)

type Logs struct {
	ID            int             `json:"id"`
	Time          sql.NullString  `json:"time"`
	RemoteAddr    sql.NullString  `json:"remote_addr"`
	XForwardedFor sql.NullString  `json:"x_forwarded_for"`
	RequestID     sql.NullString  `json:"request_id"`
	RemoteUser    sql.NullString  `json:"remote_user"`
	BytesSent     sql.NullInt32   `json:"bytes_sent"`
	RequestTime   sql.NullFloat64 `json:"request_time"`
	Status        sql.NullInt32   `json:"status"`
	Vhost         sql.NullString  `json:"vhost"`
	RequestProto  sql.NullString  `json:"request_proto"`
	Path          sql.NullString  `json:"path"`
	RequestQuery  sql.NullString  `json:"request_query"`
	RequestLength sql.NullInt32   `json:"request_length"`
	Duration      sql.NullFloat64 `json:"duration"`
	Method        sql.NullString  `json:"method"`
	HTTPReferrer  sql.NullString  `json:"http_referrer"`
	HTTPUserAgent sql.NullString  `json:"http_user_agent"`
}
type LogsJSON struct {
	ID            int     `json:"id"`
	Time          string  `json:"time"`
	RemoteAddr    string  `json:"remote_addr"`
	XForwardedFor string  `json:"x_forwarded_for"`
	RequestID     string  `json:"request_id"`
	RemoteUser    string  `json:"remote_user"`
	BytesSent     int32   `json:"bytes_sent"`
	RequestTime   float64 `json:"request_time"`
	Status        int32   `json:"status"`
	Vhost         string  `json:"vhost"`
	RequestProto  string  `json:"request_proto"`
	Path          string  `json:"path"`
	RequestQuery  string  `json:"request_query"`
	RequestLength int32   `json:"request_length"`
	Duration      float64 `json:"duration"`
	Method        string  `json:"method"`
	HTTPReferrer  string  `json:"http_referrer"`
	HTTPUserAgent string  `json:"http_user_agent"`
}

// sqlの結果をJSONに変換
func ConvertToJSON(rows *sql.Rows) ([]LogsJSON, error) {
	logEntry := Logs{}
	var result []LogsJSON
	for rows.Next() {
		// sqlの結果を取得
		err := rows.Scan(
			&logEntry.ID, &logEntry.Time, &logEntry.RemoteAddr, &logEntry.XForwardedFor, &logEntry.RequestID,
			&logEntry.RemoteUser, &logEntry.BytesSent, &logEntry.RequestTime, &logEntry.Status, &logEntry.Vhost,
			&logEntry.RequestProto, &logEntry.Path, &logEntry.RequestQuery, &logEntry.RequestLength, &logEntry.Duration,
			&logEntry.Method, &logEntry.HTTPReferrer, &logEntry.HTTPUserAgent,
		)
		if err != nil {
			return nil, err
		} else {
			// 構造体Logsを構造体LogsJSONに屁感
			logJSON := convertToType(logEntry)
			result = append(result, logJSON)
		}
	}
	return result, nil
}

// sql.NullString型をstring型に変換
func convertToType(log Logs) LogsJSON {
	return LogsJSON{
		ID:            log.ID,
		Time:          getStringFromNullString(log.Time),
		RemoteAddr:    getStringFromNullString(log.RemoteAddr),
		XForwardedFor: getStringFromNullString(log.XForwardedFor),
		RequestID:     getStringFromNullString(log.RequestID),
		RemoteUser:    getStringFromNullString(log.RemoteUser),
		BytesSent:     getIntFromNullInt(log.BytesSent),
		RequestTime:   getFloatFromNullFloat(log.RequestTime),
		Status:        getIntFromNullInt(log.Status),
		Vhost:         getStringFromNullString(log.Vhost),
		RequestProto:  getStringFromNullString(log.RequestProto),
		Path:          getStringFromNullString(log.Path),
		RequestQuery:  getStringFromNullString(log.RequestQuery),
		RequestLength: getIntFromNullInt(log.RequestLength),
		Duration:      getFloatFromNullFloat(log.Duration),
		Method:        getStringFromNullString(log.Method),
		HTTPReferrer:  getStringFromNullString(log.HTTPReferrer),
		HTTPUserAgent: getStringFromNullString(log.HTTPUserAgent),
	}
}

func getStringFromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getFloatFromNullFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0.0
}

func getIntFromNullInt(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}
