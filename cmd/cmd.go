package cmd

import (
	"github.com/darwinia-network/token/config"
	middlewares "github.com/darwinia-network/token/middleware"
	serverHttp "github.com/darwinia-network/token/server/routers/http"
	"github.com/darwinia-network/token/util"
	"github.com/urfave/cli/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func Run(args []string)  error   {

	app := &cli.App {
		Name: "Darwinia-token",
		Version: "0.1",
		Before: func(context *cli.Context) error {
			runtime.GOMAXPROCS(runtime.NumCPU())
			config.LoadConf()
			return nil

		},
		Action: func(c *cli.Context) error {
			util.GraceShutdown(&http.Server{Addr: config.Cfg.ServerHost, Handler: setupRouter()})
			return nil
		},
	}
	return app.Run(args)

}

func setupRouter() (server *gin.Engine) {
	server = gin.Default()
	server.Use(middlewares.CORS())
	serverHttp.Run(server)
	return
}


