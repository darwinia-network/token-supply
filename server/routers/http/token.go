package http

import (
	"github.com/darwinia-network/token/services/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ringSupply() gin.HandlerFunc {
	return func(c *gin.Context) {
		supply := token.RingSupply()
		if c.Query("t") == "totalSupply" {
			c.String(http.StatusOK, supply.TotalSupply.String())
			return
		}
		if c.Query("t") == "CirculatingSupply" {
			c.String(http.StatusOK, supply.CirculatingSupply.String())
			return
		}
		c.JSON(http.StatusOK, JsonFormat(supply, 0))
	}
}
func ktonSupply() gin.HandlerFunc {
	return func(c *gin.Context) {
		supply := token.KtonSupply()
		if c.Query("t") == "totalSupply" {
			c.String(http.StatusOK, supply.TotalSupply.String())
			return
		}
		if c.Query("t") == "CirculatingSupply" {
			c.String(http.StatusOK, supply.CirculatingSupply.String())
			return
		}
		c.JSON(http.StatusOK, JsonFormat(supply, 0))
	}
}
