package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/db"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/server"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var useCwd *bool

func init() {
	useCwd = flag.Bool("usecwd", false, "")
	flag.Parse()
}

func main() {
	json, err := readConnDatasJson()
	if err != nil {
		fmt.Println("An error occurred.")
		fmt.Println(err)
		os.Exit(1)
	}
	connDatas, err := db.JsonToDBConnectionData(json)
	if err != nil {
		fmt.Println("An error occurred.")
		fmt.Println(err)
		os.Exit(1)
	}

	dbMap, _ := db.CreateDBMap(connDatas)
	dbFuncDataMap := make(types.DBFuncDataMap)

	app := server.CreateServer(dbMap, dbFuncDataMap, db.RunQueryChan, db.RunExecChan, db.RowToJson, db.ResultObjectToJson)
	app.Run("localhost:3000")
}

/*
connDatas.env.json 읽기
*/
func readConnDatasJson() ([]byte, error) {
	var jsonFilePath string
	if *useCwd {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		jsonFilePath = filepath.Join(cwd, "connDatas.env.json")
	} else {
		execPath, err := os.Executable()
		if err != nil {
			return nil, err
		}

		execDir := filepath.Dir(execPath)
		jsonFilePath = filepath.Join(execDir, "connDatas.env.json")
	}

	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	return jsonContent, nil
}
