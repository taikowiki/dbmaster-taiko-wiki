package dbfunc

import (
	"context"

	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var DBFuncMap = types.DBFuncMap{
	"song": func(
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

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, "SELECT * FROM `song` WHERE `songNo` = ?", params["songNo"])
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
}
