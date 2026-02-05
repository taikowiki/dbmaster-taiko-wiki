package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taikowiki/dbmaster-taiko-wiki/src/types"
)

func permissionCheckMiddleware(dbMap types.DBMap, runQueryChan types.RunQueryChanFuncType) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		dbMasterDB := dbMap["dbMaster"]
		if dbMasterDB == nil {
			c.AbortWithStatusJSON(500, gin.H{
				"reason": "NO_DBMASTER_DATABASE",
			})
			return
		}
		apiKey := c.Request.Header.Get("X-Api-Key")
		if apiKey == "" {
			c.AbortWithStatus(401)
			return
		}

		if path == "/query" || path == "/exec" || path == "/func" {
			ch := make(chan types.RowOrError)
			go runQueryChan(ch, c, dbMasterDB, "SELECT COUNT(*) AS `count` FROM `key/all` WHERE `key` = ?", apiKey)

			err := <-ch
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}

			row := <-ch
			v, ok := row.(map[string]any)
			if !ok {
				log.Println("Permission check error.")
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}
			count, ok := v["count"].(int64)
			if !ok {
				log.Println("Permission check error.")
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}

			if count > 0 {
				c.Next()
				return
			}
		}

		if path == "/func" {
			var req FuncRequestData
			if err := c.ShouldBindJSON(&req); err != nil {
				c.AbortWithStatus(400)
				return
			}

			ch := make(chan types.RowOrError)
			go runQueryChan(ch, c, dbMasterDB, "SELECT COUNT(*) AS `count` FROM `key/func` WHERE `key` = ? AND `func` = ?", apiKey, req.Name)

			err := <-ch
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}

			row := <-ch
			v, ok := row.(map[string]any)
			if !ok {
				log.Println("Permission check error.")
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}
			count, ok := v["count"].(int64)
			if !ok {
				log.Println("Permission check error.")
				c.AbortWithStatus(500)
				log.Println(err)
				return
			}

			if count > 0 {
				c.Set("funcRequestData", req)
				c.Next()
				return
			}
		}

		c.AbortWithStatus(403)
	}
}
