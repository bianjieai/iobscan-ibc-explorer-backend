package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TokenController struct {
}

func (ctl *TokenController) PopularSymbols(c *gin.Context) {
	var req vo.PopularSymbolsReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailBadRequest(err.Error()))
		return
	}

	minHops, err := strconv.ParseInt(req.MinHops, 10, 64)
	if err != nil || minHops < 0 {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("invalid min_hops value"))
		return
	}
	minReceiveTxs, err := strconv.ParseInt(req.MinReceiveTxs, 10, 64)
	if err != nil || minReceiveTxs < 0 {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("invalid min_receive_txs value"))
		return
	}

	res, e := tokenService.PopularSymbols(int(minHops), minReceiveTxs)
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(res))
}
