package main

import (
	"fmt"
	"log"
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

	dbMap, err := db.CreateDBMap(connDatas)
	if err != nil {
		log.Panicln(err)
	}
	env, err := util.LoadEnv()
	if err != nil {
		log.Panicln(err)
	}

	app := server.CreateServer(dbMap, dbfunc.DBFuncMap, db.RunQueryChan, db.RunExecChan, db.RowToJson, db.ResultObjectToJson)
	app.Run(":" + env["PORT"])
}
