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
type RowToJsonFuncType = func(row map[string]any) ([]byte, error)

// server
type RowOrError any
