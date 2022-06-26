package api

import (
	"net/http"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/middleware"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/rest"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	store = persistence.NewInMemoryStore(time.Second)
)

func Routers(Router *gin.Engine) {
	Router.Use(middleware.Cors())
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, global.Config.App.Version)
	})

	ibcRouter := Router.Group("ibc")
	tokenPage(ibcRouter)
	channelPage(ibcRouter)
	chainPage(ibcRouter)
	relayerPage(ibcRouter)
}

func tokenPage(r *gin.RouterGroup) {
	ctl := rest.TokenController{}
	r.GET("/tokenList", ctl.List)
}

func channelPage(r *gin.RouterGroup) {

}

func chainPage(r *gin.RouterGroup) {
	ctl := rest.ChainController{}
	r.GET("/chainList", cache.CachePage(store, 3*time.Second, ctl.List))
}

func relayerPage(r *gin.RouterGroup) {
	ctl := rest.RelayerController{}
	r.GET("/relayerList", cache.CachePage(store, 3*time.Second, ctl.List))
}
