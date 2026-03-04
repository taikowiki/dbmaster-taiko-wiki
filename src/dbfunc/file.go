package dbfunc

import (
	"context"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var file_dbFuncMap = types.DBFuncMap{
	"file.getByFileName": func(
		ch chan types.ResponseJsonOrError,
		c context.Context,
		dbMap types.DBMap,
		params map[string]any,
		runQueryChan types.RunQueryChanFuncType,
		runExecChan types.RunExecChanFuncType,
		rowToJson types.RowToJsonFuncType,
		resultObjectToJson types.ResultObjectToJsonFuncType,
	) {
		defer close(ch)
		taikowikiDB := dbMap["taikowiki"]
		if taikowikiDB == nil {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		fileName, ok := params["fileName"].(string)
		if !ok || fileName == "" {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, "SELECT * FROM `file/log` WHERE `fileName` = ?", fileName)
		err := <-receiveChan
		if err != nil {
			ch <- err
			return
		}
		ch <- nil

		for row := range receiveChan {
			if v, ok := row.(map[string]any); ok {
				json, err := rowToJson(v)
				if err != nil {
					continue
				}
				ch <- json
			}
		}
	},
	"file.newLog": func(
		ch chan types.ResponseJsonOrError,
		c context.Context,
		dbMap types.DBMap,
		params map[string]any,
		runQueryChan types.RunQueryChanFuncType,
		runExecChan types.RunExecChanFuncType,
		rowToJson types.RowToJsonFuncType,
		resultObjectToJson types.ResultObjectToJsonFuncType,
	) {
		defer close(ch)
		taikowikiDB := dbMap["taikowiki"]
		if taikowikiDB == nil {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		uuid, ok1 := params["UUID"].(string)
		originalFileName, ok2 := params["originalFileName"].(string)
		fileName, ok3 := params["fileName"].(string)

		if !ok1 || !ok2 || !ok3 || uuid == "" || originalFileName == "" || fileName == "" {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.ResultObjectOrError)
		go runExecChan(receiveChan, c, taikowikiDB, "INSERT INTO `file/log` (`UUID`, `originalFileName`, `fileName`) VALUES (?, ?, ?)", uuid, originalFileName, fileName)
		err := <-receiveChan
		if err != nil {
			ch <- err
			return
		}
		ch <- nil

		result := <-receiveChan
		if v, ok := result.(types.ResultObject); ok {
			json, err := resultObjectToJson(v)
			if err == nil {
				ch <- json
				return
			}
		}
	},
}
