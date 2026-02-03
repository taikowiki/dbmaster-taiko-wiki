package dbfunc

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var rating_dbFuncMap = types.DBFuncMap{
	/*
		@types `params.uuid` string
	*/
	"rating.simple-profile": func(
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
		ratingDB := dbMap["rating"]
		if ratingDB == nil {
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

		query := "SELECT `user/taiko_profile`.`taikoNumber`, `user/taiko_profile`.`UUID`, `user/taiko_profile`.`nickname`, `user/rating_data`.`currentRatingScore`, `user/rating_data`.`lastUpload` FROM `user/taiko_profile` "
		query += "JOIN `user/rating_data` ON `user/taiko_profile`.`UUID` = `user/rating_data`.`UUID` "
		query += "WHERE `user/taiko_profile`.`UUID` = ?"
		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, ratingDB, query, uuid)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
	/*
		@types `params.uuids` []string
	*/
	"rating.taiko-profiles": func(
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
		ratingDB := dbMap["rating"]
		if ratingDB == nil {
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
		query, _, err := goqu.Dialect("mysql").From("user/taiko_profile").Select("taikoNumber", "nickname", "UUID").Where(goqu.Ex{
			"UUID": uuids,
		}).ToSQL()
		if err != nil {
			ch <- &types.ErrorWithStatus{
				Status: 400,
			}
			return
		}

		receiveChan := make(chan types.RowOrError)
		go runQueryChan(receiveChan, c, ratingDB, query)
		ch <- checkReceiveChanError(receiveChan)

		sendRowJson(receiveChan, ch, rowToJson, resultObjectToJson)
	},
}
