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
	taskTools(ibcRouter)
	addressPage(ibcRouter)
	overviewPage(ibcRouter)
}

func homePage(r *gin.RouterGroup) {
	ctl := rest.HomeController{}
	r.GET("/chains", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.DailyChains))
	r.GET("/chains_connection", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.ChainsConnection))
	r.GET("/baseDenoms", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.AuthDenoms))
	r.GET("/denoms", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IbcDenoms))
	r.GET("/statistics", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.Statistics))
	r.POST("/searchPoint", ctl.SearchPoint)
}

func txsPage(r *gin.RouterGroup) {
	ctl := rest.IbcTransferController{}
	r.GET("/txs", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTxs))
	r.GET("/txs_detail/:hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTxDetailNew))
	r.GET("/trace_source/:hash", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TraceSource))
	r.GET("/txs/searchCondition", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.SearchCondition))
}

func tokenPage(r *gin.RouterGroup) {
	ctl := rest.TokenController{}
	r.GET("/tokenList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.List))
	r.GET("/ibcTokenList", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.IBCTokenList))
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
	r.GET("/relayer/:relayer_id", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.Detail))
	r.GET("/relayer/:relayer_id/txs", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.DetailRelayerTxs))
	r.GET("/relayer/names", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.RelayerNameList))
	r.GET("/relayer/:relayer_id/relayedTrend", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.RelayerTrend))
	r.POST("/relayerCollect", ctl.Collect)
	r.GET("/relayer/:relayer_id/transferTypeTxs", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TransferTypeTxs))
	r.GET("/relayer/:relayer_id/totalRelayedValue", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TotalRelayedValue))
	r.GET("/relayer/:relayer_id/totalFeeCost", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TotalFeeCost))
}

func addressPage(r *gin.RouterGroup) {
	ctl := rest.AddressController{}
	r.GET("/chain/:chain/address/:address", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.BaseInfo))
	r.GET("/chain/:chain/address/:address/tokens", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TokenList))
	r.GET("/chain/:chain/address/:address/txs", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TxsList))
	r.GET("/chain/:chain/address/:address/txs/export", ctl.TxsExport)
	r.GET("/chain/:chain/address/:address/accounts", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.AccountList))
}

func overviewPage(r *gin.RouterGroup) {
	ctl := rest.OverviewController{}
	r.GET("/overview/market_heatmap", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.MarketHeatmap))
	r.GET("/overview/chain_volume", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.ChainVolume))
	r.GET("/overview/chain_volume_trend", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.ChainVolumeTrend))
	r.GET("/overview/token_distribution", cache.CachePage(store, time.Duration(aliveSeconds)*time.Second, ctl.TokenDistribution))
}

func cacheTools(r *gin.RouterGroup) {
	ctl := rest.CacheController{}
	r.DELETE("/cache/:key", ctl.Del)
}

func taskTools(r *gin.RouterGroup) {
	ctl := rest.TaskController{}
	r.POST("/task/:task_name", ctl.Run)
}
