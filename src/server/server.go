package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

type QueryRequestData struct {
	Name  string `json:"name"`
	Query string `json:"query"`
	Args  []any  `json:"args"`
}

func CreateServer(dbMap map[string](*sql.DB), runQueryChan types.RunQueryChanFuncType, rowToJson types.RowToJsonFuncType) *gin.Engine {
	r := gin.Default()
	r.POST("/query", CreateQueryHandler(dbMap["taikowiki"], runQueryChan, rowToJson))
	return r
}

func CreateQueryHandler(
	db *sql.DB,
	runQueryChan types.RunQueryChanFuncType,
	rowToJson types.RowToJsonFuncType,
) func(*gin.Context) {
	return func(c *gin.Context) {
		var req QueryRequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Status(400)
			return
		}

		ch := make(chan types.RowOrError)
		go runQueryChan(ch, c, db, "SELECT * FROM `song` LIMIT 5")

		// error check
		err := <-ch
		if err == nil {
			c.Writer.Header().Set("Content-Type", "application/x-ndjson")
			c.Writer.WriteHeader(200)
		} else {
			c.Status(404)
			return
		}

		flusher := c.Writer.(http.Flusher)
		for row := range ch {
			if v, ok := row.(map[string]any); ok {
				json, err := rowToJson(v)
				if err != nil {
					continue
				}
				fmt.Fprintln(c.Writer, (string)(json))
				flusher.Flush()
			}
		}
	}
}
