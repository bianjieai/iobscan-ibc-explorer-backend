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
	homePage(ibcRouter)
	txsPage(ibcRouter)
	tokenPage(ibcRouter)
	channelPage(ibcRouter)
	chainPage(ibcRouter)
	relayerPage(ibcRouter)
	cacheTools(ibcRouter)
}

func homePage(r *gin.RouterGroup) {
	ctl := rest.HomeController{}
	r.GET("/chains", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.DailyChains))
	r.GET("/baseDenoms", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcBaseDenoms))
	r.GET("/denoms", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcDenoms))
	r.GET("/statistics", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.Statistics))
	r.POST("/searchPoint", ctl.SearchPoint)
}

func txsPage(r *gin.RouterGroup) {
	ctl := rest.IbcTransferController{}
	r.GET("/txs", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTxs))
	r.GET("/txs/:hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTxDetail))
	r.GET("/txs_detail/:hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTxDetailNew))
	r.GET("/trace_source/:hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TraceSource))
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
	r.POST("/relayerCollect", ctl.Collect)
}

func cacheTools(r *gin.RouterGroup) {
	ctl := rest.CacheController{}
	r.DELETE("/:key", ctl.Del)
}
