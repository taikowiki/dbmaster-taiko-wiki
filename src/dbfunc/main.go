package dbfunc

import (
	"maps"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

var DBFuncMap = types.DBFuncMap{}

func InitDBFuncMap() {
	maps.Copy(DBFuncMap, song_dbFuncMap)
	maps.Copy(DBFuncMap, user_dbFuncMap)
	maps.Copy(DBFuncMap, rating_dbFuncMap)
}

func checkReceiveChanError(receiveChan any) error {
	var err any
	switch ch := receiveChan.(type) {
	case chan types.RowOrError:
		{
			err = <-ch
		}
	case chan types.ResultObjectOrError:
		{
			err = <-ch
		}
	}

	if err != nil {
		return err.(error)
	} else {
		return nil
	}
}
func sendRowJson(receiveChan chan types.RowOrError, sendChan chan types.ResponseJsonOrError, rowToJson types.RowToJsonFuncType, resultObjectToJson types.ResultObjectToJsonFuncType) {
	for data := range receiveChan {
		switch v := data.(type) {
		case map[string]any:
			{
				json, err := rowToJson(v)
				if err != nil {
					continue
				}
				sendChan <- json
			}
		}
	}
}
func sendResultObjectJson(receiveChan chan types.ResultObjectOrError, sendChan chan types.ResponseJsonOrError, rowToJson types.RowToJsonFuncType, resultObjectToJson types.ResultObjectToJsonFuncType) {
	for data := range receiveChan {
		switch v := data.(type) {
		case types.ResultObject:
			{
				json, err := resultObjectToJson(v)
				if err != nil {
					continue
				}
				sendChan <- json
			}
		}
	}
}
