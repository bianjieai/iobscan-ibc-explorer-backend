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
// @Tags api_support
// @Accept  json
// @Produce  json
// @Success 200 {object} vo.StatisticInfoResp	"success"
// @Router /data/statistics/api_support [get]
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
// @Tags api_support
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Success 200 {object} vo.FailTxsListResp	"success"
// @Router /data/fail_txs/api_support [get]
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
// @Tags api_support
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Param   tx_hash     query   string    false "tx_hash"
// @Param   chain_id     query   string    false "chain_id"
// @Success 200 {object} vo.RelayerTxFeesResp	"success"
// @Router /data/relayers_fee/api_support [get]
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
// @Tags api_support
// @Accept  json
// @Produce  json
// @Success 200 {object} vo.AccountsDailyResp	"success"
// @Router /data/accounts_daily/api_support [get]
func (ctl *ApiSupportController) AccountsDaily(c *gin.Context) {
	resp, err := staticInfoService.AccountsDailyStatistic()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

// @Summary list
// @Description get IBC BUSD of chains
// @ID list
// @Tags api_support
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Param   use_count   query   bool false    "if used count" Enums(true, false)
// @Success 200 {object} vo.ChainListResp	"success"
// @Router /data/chainList/api_support [get]
func (ctl *ApiSupportController) List(c *gin.Context) {
	var req vo.ChainListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	if req.UseCount {
		total, err := chainService.Count()
		if err != nil {
			c.JSON(http.StatusOK, response.FailError(err))
			return
		}
		c.JSON(http.StatusOK, response.Success(total))
		return
	}
	resp, err := chainService.List(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
