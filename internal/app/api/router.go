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

	//Router.Use(middleware.Logger())

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, global.Config.App.Version)
	})

	ibcRouter := Router.Group("ibc")
	txCtl(ibcRouter)
	chainCtl(ibcRouter)
	taskTools(ibcRouter)
	feeCtl(ibcRouter)
	addressCtl(ibcRouter)
	tokenCtl(ibcRouter)
}

func txCtl(r *gin.RouterGroup) {
	ctl := rest.IbcTxController{}
	r.GET("/txs/:tx_hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.Query))
	r.GET("/transfers/statistics/:chain/failure", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.FailureStatistics))
	r.GET("/transfers/statistics/:chain/flow", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.FlowInfoStatistics))
}

func chainCtl(r *gin.RouterGroup) {
	ctl := rest.ChainController{}
	r.GET("/chains", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
	r.GET("/chains/statistics", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcChainsNum))
	r.GET("/chains/volume", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcChainsVolume))
	r.GET("/chains/active", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcChainsActive))

}

func taskTools(r *gin.RouterGroup) {
	ctl := rest.TaskController{}
	r.POST("/task/:task_name", ctl.Run)
}

func feeCtl(r *gin.RouterGroup) {
	ctl := rest.IbcFeeController{}
	r.GET("/fee/statistics/:chain/paid", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.ChinFeeStatistics))
}

func addressCtl(r *gin.RouterGroup) {
	ctl := rest.AddressController{}
	r.GET("/addresses/statistics/:chain", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.ChainAddressStatistics))
	r.GET("/addresses/statistics", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.AllChainAddressStatistics))
}

func tokenCtl(r *gin.RouterGroup) {
	ctl := rest.TokenController{}
	r.GET("/tokens/popular-symbols", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.PopularSymbols))
}
