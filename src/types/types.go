package types

import (
	"context"
	"database/sql"
)

// db
type DBConnectionData struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"` // Assuming JSON key is "db_name"
	Timezone string `json:"timezone"`
}

type RunQueryChanFuncType = func(ch chan RowOrError, ctx context.Context, db *sql.DB, query string, args ...any)
type RunExecChanFuncType = func(ch chan ResultObjectOrError, ctx context.Context, db *sql.DB, query string, args ...any)
type RowToJsonFuncType = func(row map[string]any) ([]byte, error)
type ResultObjectToJsonFuncType = func(resultObject ResultObject) ([]byte, error)

type DBMap = map[string](*sql.DB)
type ResultObject struct {
	LastInsertId int64 `json:"lastInsertId"`
	RowsAffected int64 `json:"rowsAffected"`
}

// server
type RowOrError any
type ResultObjectOrError any
