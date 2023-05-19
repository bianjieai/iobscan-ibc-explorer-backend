package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChainController struct {
}

func (ctl *ChainController) List(c *gin.Context) {
	list, e := chainService.List()
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(list))
}

func (ctl *ChainController) IbcChainsNum(c *gin.Context) {
	resp, e := chainService.IbcChainsNum()
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *ChainController) IbcChainsVolume(c *gin.Context) {
	var req vo.IbcChainsVolumeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailBadRequest(err.Error()))
		return
	}

	resp, e := chainService.IbcChainsVolume(req.Chain)
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *ChainController) IbcChainsActive(c *gin.Context) {
	resp, e := chainService.IbcChainsActive()
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}
