package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
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

func (ctl *ChainController) ActiveChainNum(c *gin.Context) {
	resp, e := chainService.ActiveChainNum()
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}
