package dbfunc

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var user_dbFuncMap = types.DBFuncMap{
	/*
		@type `params.provider` string
		@type `params.providerId` string
	*/
	"user.user-data": func(
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

		provider, ok1 := params["provider"].(string)
		providerId, ok2 := params["providerId"].(string)
		if !ok1 || !ok2 {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, "SELECT * FROM `user/data` WHERE `provider` = ? AND `providerId` = ?", provider, providerId)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
	/*
		@type `params.uuid` string
	*/
	"user.user-data-by-uuid": func(
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

		uuid, ok := params["uuid"].(string)
		if !ok {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, taikowikiDB, "SELECT * FROM `user/data` WHERE `UUID` = ?", uuid)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
	"user.nick-and-uuid": func(
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
		go runQueryChan(receiveChan, c, taikowikiDB, "SELECT `nickname`, `UUID` FROM `user/data`")
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
	/*
		@types `params.uuids` []string
	*/
	"user.nicks": func(
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

		uuids, ok := params["uuids"].([]any)
		if !ok {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		for _, v := range uuids {
			if _, ok := v.(string); !ok {
				ch <- &types.ErrorWithStatus{
					Status: 400,
				}
				return
			}
		}

		var err any
		query, _, err := goqu.Dialect("mysql").From("user/data").Select("nickname", "UUID").Where(goqu.Ex{
			"UUID": uuids,
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
}
