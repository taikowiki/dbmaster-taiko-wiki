package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

/*
연결을 생성
*/
func CreateDB(connData types.DBConnectionData) (*sql.DB, error) {
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

func CreateDBMap(connDatas []types.DBConnectionData) (types.DBMap, error) {
	connectionMap := make(types.DBMap)
	for _, data := range connDatas {
		connection, err := CreateDB(data)
		if err != nil {
			return nil, err
		}
		connectionMap[data.Name] = connection
	}
	return connectionMap, nil
}

/*
json을 DBConnectionData로 변환
*/
func JsonToDBConnectionData(jsonContent []byte) ([]types.DBConnectionData, error) {
	var connDatas []types.DBConnectionData
	err := json.Unmarshal(jsonContent, &connDatas)
	if err != nil {
		return nil, err
	}
	return connDatas, nil
}

/*
Query를 실행한 후 row를 채널로 보냄
*/
func RunQueryChan(ch chan types.RowOrError, ctx context.Context, db *sql.DB, query string, args ...any) {
	defer close(ch)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		ch <- err
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		ch <- err
		return
	}
	if len(columns) == 0 {
		ch <- err
		return
	}

	// no error
	ch <- nil

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
			return
		}
	}
}

/*
Row를 json으로 변환
*/
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

/*
Exec를 실행한 후 ResultObject를 채널로 보냄
*/
func RunExecChan(ch chan types.ResultObjectOrError, ctx context.Context, db *sql.DB, query string, args ...any) {
	defer close(ch)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		ch <- err
		return
	}

	// no error
	ch <- nil

	resultObject := types.ResultObject{
		LastInsertId: 0,
		RowsAffected: 0,
	}

	lastInsertedId, err := result.LastInsertId()
	if err == nil {
		resultObject.LastInsertId = lastInsertedId
	}
	rowsAffected, err := result.RowsAffected()
	if err == nil {
		resultObject.RowsAffected = rowsAffected
	}

	ch <- resultObject
}

/*
ResultObject를 json으로 변환
*/
func ResultObjectToJson(resultObject types.ResultObject) ([]byte, error) {
	return json.Marshal(resultObject)
}
