package main

import (
	"fmt"
	"os"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/db"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/dbfunc"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/server"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/util"
)

func init() {
	util.LoadFlag()
	dbfunc.InitDBFuncMap()
}

func main() {
	connDatas, err := util.LoadConnDatas()
	if err != nil {
		fmt.Println("An error occurred.")
		fmt.Println(err)
		os.Exit(1)
	}

	dbMap, _ := db.CreateDBMap(connDatas)

	app := server.CreateServer(dbMap, dbfunc.DBFuncMap, db.RunQueryChan, db.RunExecChan, db.RowToJson, db.ResultObjectToJson)
	app.Run("localhost:3000")
}
