package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiSupportController struct {
}

func (ctl *ApiSupportController) StatisticInfo(c *gin.Context) {
	resp, err := staticInfoService.IbcTxStatistic()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *ApiSupportController) FailTxsList(c *gin.Context) {
	var req vo.FailTxsListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	resp, err := ibcTxService.ListFailTxs(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *ApiSupportController) RelayerTxsFee(c *gin.Context) {
	var req vo.RelayerTxFeesReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	resp, err := ibcTxService.ListRelayerTxFees(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
