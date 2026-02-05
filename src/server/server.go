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
type FuncRequestData struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
}

/*
Gin 서버 생성
*/
func CreateServer(
	dbMap types.DBMap,
	dbFuncMap types.DBFuncMap,
	runQueryChan types.RunQueryChanFuncType,
	runExecChan types.RunExecChanFuncType,
	rowToJson types.RowToJsonFuncType,
	resultObjectToJson types.ResultObjectToJsonFuncType,
) *gin.Engine {
	r := gin.Default()
	r.Use(permissionCheckMiddleware(dbMap, runQueryChan))
	r.POST("/query", createQueryHandler(dbMap, runQueryChan, rowToJson))
	r.POST("/exec", createExecHandler(dbMap, runExecChan, resultObjectToJson))
	r.POST("/func", createDBFuncHandler(dbMap, dbFuncMap, runQueryChan, runExecChan, rowToJson, resultObjectToJson))
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

		sendQueryRows(c, db, req.Query, req.Args, runQueryChan, rowToJson)
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

		sendExecResult(c, db, req.Query, req.Args, runExecChan, resultObjectToJson)
	}
}

func createDBFuncHandler(
	dbMap types.DBMap,
	dbFuncMap types.DBFuncMap,
	runQueryChan types.RunQueryChanFuncType,
	runExecChan types.RunExecChanFuncType,
	rowToJson types.RowToJsonFuncType,
	resultObjectToJson types.ResultObjectToJsonFuncType,
) func(*gin.Context) {
	return func(c *gin.Context) {
		var req FuncRequestData
		funcRequestData, exists := c.Get("funcRequestData")
		if exists {
			req = funcRequestData.(FuncRequestData)
		} else {
			if err := c.ShouldBindJSON(&req); err != nil {
				c.Status(400)
				return
			}
		}

		dbFunc := dbFuncMap[req.Name]
		if dbFunc == nil {
			c.Status(400)
			return
		}

		ch := make(chan types.ResponseJsonOrError)
		go dbFunc(ch, c, dbMap, req.Params, runQueryChan, runExecChan, rowToJson, resultObjectToJson)
		err := <-ch
		if err == nil {
			c.Writer.Header().Set("Content-Type", "application/x-ndjson")
			c.Writer.WriteHeader(200)
		} else if e, ok := err.(*types.ErrorWithStatus); ok {
			c.Status(e.Status)
			return
		} else {
			c.Status(500)
			return
		}

		flusher := c.Writer.(http.Flusher)
		for row := range ch {
			if v, ok := row.([]byte); ok {
				fmt.Fprintln(c.Writer, (string)(v))
				flusher.Flush()
			}
		}
	}
}

func sendQueryRows(
	c *gin.Context,
	db *sql.DB,
	query string,
	args []any,
	runQueryChan types.RunQueryChanFuncType,
	rowToJson types.RowToJsonFuncType,
) {
	ch := make(chan types.RowOrError)
	go runQueryChan(ch, c, db, query, args...)

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

func sendExecResult(
	c *gin.Context,
	db *sql.DB,
	query string,
	args []any,
	runExecChan types.RunExecChanFuncType,
	resultObjectToJson types.ResultObjectToJsonFuncType,
) {
	ch := make(chan types.ResultObjectOrError)
	go runExecChan(ch, c, db, query, args...)

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
