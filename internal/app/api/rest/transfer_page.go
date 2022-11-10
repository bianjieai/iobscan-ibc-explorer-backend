package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IbcTransferController struct {
}

func (ctl *IbcTransferController) TransferTxs(c *gin.Context) {
	var req vo.TranaferTxsReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}
	if req.UseCount {
		count, err := transferService.TransferTxsCount(&req)
		if err != nil {
			c.JSON(http.StatusOK, response.FailError(err))
			return
		}
		c.JSON(http.StatusOK, response.Success(count))
		return
	}
	resp, err := transferService.TransferTxs(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *IbcTransferController) TransferTxDetail(c *gin.Context) {
	hash := c.Param("hash")
	resp, err := transferService.TransferTxDetail(hash)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *IbcTransferController) TransferTxDetailNew(c *gin.Context) {
	hash := c.Param("hash")
	resp, err := transferService.TransferTxDetailNew(hash)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *IbcTransferController) TraceSource(c *gin.Context) {
	hash := c.Param("hash")
	var req vo.TraceSourceReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}
	resp, err := transferService.TraceSource(hash, &req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *IbcTransferController) SearchCondition(c *gin.Context) {
	resp, err := transferService.SearchCondition()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
