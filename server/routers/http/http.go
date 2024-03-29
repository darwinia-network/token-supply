package http

import (
	middlewares "github.com/darwinia-network/token/middleware"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"time"
)

func Run(server *gin.Engine) {

	store := persistence.NewInMemoryStore(time.Minute)
	api := server.Group("/api")

	server.GET("supply/ring", middlewares.PageCache(store, time.Second*10, ringSupply()))
	server.GET("supply/kton", middlewares.PageCache(store, time.Second*10, ktonSupply()))
	api.GET("/status", func(c *gin.Context) {
		c.String(200, "OK")
	})
	api.GET("supply", middlewares.PageCache(store, time.Second*10, ringSupply()))

}

func JsonFormat(data interface{}, code int) map[string]interface{} {
	r := gin.H{
		"data": data,
		"code": code,
		"msg":  responseCode[code],
	}
	return r
}

var responseCode = map[int]string{
	0:    "ok",
	1001: "params error",
}
