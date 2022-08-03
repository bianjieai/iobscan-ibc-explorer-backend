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
	store        = persistence.NewInMemoryStore(time.Second)
	aliveSeconds = 3
)

func SetApiCacheAliveTime(duration int) {
	aliveSeconds = duration
}

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
	cacheTools(ibcRouter)
}

func tokenPage(r *gin.RouterGroup) {
	ctl := rest.TokenController{}
	r.GET("/tokenList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
	r.GET("/:base_denom/ibcTokenList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IBCTokenList))
}

func channelPage(r *gin.RouterGroup) {
	ctl := rest.ChannelController{}
	r.GET("/channelList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
}

func chainPage(r *gin.RouterGroup) {
	ctl := rest.ChainController{}
	r.GET("/chainList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
}

func relayerPage(r *gin.RouterGroup) {
	ctl := rest.RelayerController{}
	r.GET("/relayerList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
}

func cacheTools(r *gin.RouterGroup) {
	ctl := rest.CacheController{}
	r.DELETE("/:key", ctl.Del)
}
