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

// @Summary list
// @Description get ibc txs
// @ID list
// @Tags app_version
// @Accept  json
// @Produce  json
// @Success 200 {object} vo.StatisticInfoResp	"success"
// @Router /ibc/statistics/api_support [get]
func (ctl *ApiSupportController) StatisticInfo(c *gin.Context) {
	resp, err := staticInfoService.IbcTxStatistic()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

// @Summary list
// @Description get  fail txs list
// @ID list
// @Tags app_version
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Success 200 {object} vo.FailTxsListResp	"success"
// @Router /ibc/fail_txs/api_support [get]
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

// @Summary list
// @Description get relayers fee list of chains
// @ID list
// @Tags app_version
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Success 200 {object} vo.RelayerTxFeesResp	"success"
// @Router /ibc/relayers_fee/api_support [get]
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

// @Summary list
// @Description get daily accounts of chains
// @ID list
// @Tags app_version
// @Accept  json
// @Produce  json
// @Success 200 {object} vo.AccountsDailyResp	"success"
// @Router /ibc/accounts_daily/api_support [get]
func (ctl *ApiSupportController) AccountsDaily(c *gin.Context) {
	resp, err := staticInfoService.AccountsDailyStatistic()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
