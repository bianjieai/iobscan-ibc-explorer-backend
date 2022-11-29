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
	if global.Config.App.EnableSignature {
		Router.Use(middleware.SignatureVerification())
	}
	if global.Config.App.EnableRateLimit {
		Router.Use(middleware.RateLimit())
	}
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, global.Config.App.Version)
	})

	ibcRouter := Router.Group("ibc")
	txCtl(ibcRouter)
	chainCtl(ibcRouter)
}

func txCtl(r *gin.RouterGroup) {
	ctl := rest.IbcTxController{}
	r.GET("/txs/:tx_hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.Query))
	r.GET("/transfers/statistics/:chain/failure", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.FailureStatistics))
}

func chainCtl(r *gin.RouterGroup) {
	ctl := rest.ChainController{}
	r.GET("/chains", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
}
