package rest

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
)

type OverviewController struct {
}

func (ctl *OverviewController) MarketHeatmap(c *gin.Context) {
	resp, err := overviewService.MarketHeatmap()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *OverviewController) ChainVolume(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}

func (ctl *OverviewController) ChainVolumeTrend(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}

func (ctl *OverviewController) TokenDistribution(c *gin.Context) {

	var req vo.TokenDistributionReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	if req.BaseDenom == "" || req.BaseDenomChain == "" {
		c.JSON(http.StatusOK, response.FailBadRequest(fmt.Errorf("invalid parameters")))
		return
	}

	resp, err := overviewService.TokenDistribution(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
