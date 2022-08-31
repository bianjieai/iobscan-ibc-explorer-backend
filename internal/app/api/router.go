package api

import (
	_ "github.com/bianjieai/iobscan-ibc-explorer-backend/docs"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/middleware"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/rest"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
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

	ibcRouter := Router.Group("data")

	//api_support
	statisticApiSupport(ibcRouter)

}

func statisticApiSupport(r *gin.RouterGroup) {
	ctl := rest.ApiSupportController{}
	r.GET("/statistics/api_support", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.StatisticInfo))
	r.GET("/fail_txs/api_support", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.FailTxsList))
	r.GET("/relayers_fee/api_support", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.RelayerTxsFee))
	r.GET("/accounts_daily/api_support", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.AccountsDaily))
	r.GET("/chainList/api_support", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))

}
