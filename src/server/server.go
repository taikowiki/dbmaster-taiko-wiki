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

/*
Gin 서버 생성
*/
func CreateServer(
	dbMap types.DBMap,
	runQueryChan types.RunQueryChanFuncType,
	runExecChan types.RunExecChanFuncType,
	rowToJson types.RowToJsonFuncType,
	resultObjectToJson types.ResultObjectToJsonFuncType,
) *gin.Engine {
	r := gin.Default()
	r.POST("/query", createQueryHandler(dbMap, runQueryChan, rowToJson))
	r.POST("/exec", createExecHandler(dbMap, runExecChan, resultObjectToJson))
	return r
}

/*
Query 핸들러
*/
func createQueryHandler(
	dbMap types.DBMap,
	runQueryChan types.RunQueryChanFuncType,
	rowToJson types.RowToJsonFuncType,
) func(*gin.Context) {
	return func(c *gin.Context) {
		var req QueryRequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Status(400)
			return
		}

		var db *sql.DB = dbMap[req.Name]
		if db == nil {
			c.Status(400)
			return
		}

		ch := make(chan types.RowOrError)
		go runQueryChan(ch, c, db, req.Query, req.Args...)

		// error check
		err := <-ch
		if err == nil {
			c.Writer.Header().Set("Content-Type", "application/x-ndjson")
			c.Writer.WriteHeader(200)
		} else {
			fmt.Println(err)
			c.Status(400)
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

/*
Exec 핸들러
*/
func createExecHandler(
	dbMap types.DBMap,
	runExecChan types.RunExecChanFuncType,
	resultObjectToJson types.ResultObjectToJsonFuncType,
) func(*gin.Context) {
	return func(c *gin.Context) {
		var req QueryRequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Status(400)
			return
		}

		var db *sql.DB = dbMap[req.Name]
		if db == nil {
			c.Status(400)
			return
		}

		ch := make(chan types.ResultObjectOrError)
		go runExecChan(ch, c, db, req.Query, req.Args...)

		// error check
		err := <-ch
		if err != nil {
			fmt.Println(err)
			c.Status(400)
			return
		}

		result := <-ch
		if v, ok := result.(types.ResultObject); ok {
			json, err := resultObjectToJson(v)
			if err == nil {
				c.Writer.Header().Set("Content-Type", "application/json")
				c.Writer.WriteHeader(200)
				c.Writer.Write(json)
				return
			}
		}

		c.Status(500)
	}
}
