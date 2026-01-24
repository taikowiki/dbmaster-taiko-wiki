package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/db"
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

	dbMap, _ := db.CreateConnectionMap(connDatas)

	taikoWikiDB := dbMap["taikowiki"]
	if taikoWikiDB == nil {
		return
	}

	ch := make(chan map[string]any)

	ctx, _ := context.WithCancel(context.Background())
	go db.RunQueryChan(ch, ctx, taikoWikiDB, "SELECT * FROM `log` LIMIT 10")

	for row := range ch {
		json, _ := db.RowToJson(row)
		fmt.Println((string)(json))
	}
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
