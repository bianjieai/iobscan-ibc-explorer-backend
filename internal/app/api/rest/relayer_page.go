package rest

import (
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
)

type RelayerController struct {
}

func (ctl *RelayerController) List(c *gin.Context) {
	var req vo.RelayerListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	var res interface{}
	var err errors.Error
	if req.UseCount {
		res, err = relayerService.ListCount(&req)
	} else {
		res, err = relayerService.List(&req)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}

func (ctl *RelayerController) Collect(c *gin.Context) {
	filepath := c.PostForm("filepath")
	if err := relayerService.Collect(filepath); err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (ctl *RelayerController) Detail(c *gin.Context) {
	relayerId := c.Param("relayer_id")
	var res interface{}
	var err errors.Error
	res, err = relayerService.Detail(relayerId)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}

func (ctl *RelayerController) DetailRelayerTxs(c *gin.Context) {
	relayerId := c.Param("relayer_id")
	var req vo.DetailRelayerTxsReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	chainInfo, err := relayerService.CheckRelayerParams(relayerId, req.Chain)
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}
	req.Addresses = chainInfo.Addresses
	var res interface{}
	if req.UseCount {
		res, err = relayerService.DetailRelayerTxsCount(&req)
	} else {
		res, err = relayerService.DetailRelayerTxs(&req)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
