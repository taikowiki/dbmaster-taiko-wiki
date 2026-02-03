package dbfunc

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var song_dbFuncMap = types.DBFuncMap{
	"song.song-data": func(
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
	"song.partial-song-data": func(
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

		columns, ok := params["columns"].([]any)
		if !ok {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		var err any
		query, _, err := goqu.Dialect("mysql").From("song").Select(columns...).Where(goqu.Ex{
			"songNo": params["songNo"],
		}).ToSQL()
		if err != nil {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, query)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
	"song.partial-data-of-all-songs": func(
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

		columns, ok := params["columns"].([]any)
		if !ok {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		var err any
		query, _, err := goqu.Dialect("mysql").From("song").Select(columns...).ToSQL()
		if err != nil {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}
		query += "ORDER BY CAST(`songNo` AS INT) ASC"

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, query)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
}
