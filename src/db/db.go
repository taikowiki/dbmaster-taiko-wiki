package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConnectionData struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"` // Assuming JSON key is "db_name"
	Timezone string `json:"timezone"`
}

/*
연결을 생성
*/
func CreateConnection(connData DBConnectionData) (*sql.DB, error) {
	dsnParams := url.Values{}
	dsnParams.Set("parseTime", "true")
	dsnParams.Set("charset", "utf8mb4")
	if connData.Timezone != "" {
		dsnParams.Set("loc", connData.Timezone)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", connData.User, connData.Password, connData.Host, connData.Port, connData.Database, dsnParams.Encode())

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateConnectionMap(connDatas []DBConnectionData) (map[string](*sql.DB), error) {
	connectionMap := map[string](*sql.DB){}
	for _, data := range connDatas {
		connection, err := CreateConnection(data)
		if err != nil {
			return nil, err
		}
		connectionMap[data.Database] = connection
	}
	return connectionMap, nil
}

/*
json을 DBConnectionData로 변환
*/
func JsonToDBConnectionData(jsonContent []byte) ([]DBConnectionData, error) {
	var connDatas []DBConnectionData
	err := json.Unmarshal(jsonContent, &connDatas)
	if err != nil {
		return nil, err
	}
	return connDatas, nil
}

/*
Query를 실행한 후 row를 채널로 보냄
*/
func RunQueryChan(ch chan map[string]any, ctx context.Context, db *sql.DB, query string, args ...any) error {
	defer close(ch)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return nil
	}

	for rows.Next() {
		values := make([]any, len(columns))
		valuesPtrs := make([]any, len(columns))
		for i := range values {
			valuesPtrs[i] = &values[i]
		}

		err := rows.Scan(valuesPtrs...)
		if err != nil {
			continue
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			switch t := val.(type) {
			case []byte:
				row[col] = (string)(t)
			case sql.NullString:
				if t.Valid {
					row[col] = t.String
				} else {
					row[col] = nil
				}
			case sql.NullInt64:
				if t.Valid {
					row[col] = t.Int64
				} else {
					row[col] = nil
				}
			case sql.NullTime:
				if t.Valid {
					row[col] = t.Time
				} else {
					row[col] = nil
				}
			default:
				row[col] = val
			}
		}

		select {
		case ch <- row:
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func RowToJson(row map[string]any) ([]byte, error) {
	normalizedMap := make(map[string]any)
	for key, value := range row {
		if v, ok := value.(time.Time); ok {
			normalizedMap[key] = map[string]int64{
				"unixMilli": v.UnixMilli(),
			}
		} else {
			normalizedMap[key] = value
		}
	}

	return json.Marshal(normalizedMap)
}
