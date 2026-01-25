package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/db"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/server"
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

	app := server.CreateServer(dbMap, db.RunQueryChan, db.RowToJson)
	app.Run("localhost:3000")
}

/*
connDatas.json 읽기
*/
func readConnDatasJson() ([]byte, error) {
	var jsonFilePath string
	if *useCwd {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		jsonFilePath = filepath.Join(cwd, "connDatas.json")
	} else {
		execPath, err := os.Executable()
		if err != nil {
			return nil, err
		}

		execDir := filepath.Dir(execPath)
		jsonFilePath = filepath.Join(execDir, "connDatas.json")
	}

	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	return jsonContent, nil
}
