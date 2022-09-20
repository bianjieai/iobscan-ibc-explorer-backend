package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeController struct {
}

func (ctl *HomeController) DailyChains(c *gin.Context) {
	resp, err := homeService.DailyChains()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *HomeController) IbcBaseDenoms(c *gin.Context) {
	resp, err := homeService.IbcBaseDenoms()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *HomeController) IbcDenoms(c *gin.Context) {
	resp, err := homeService.IbcDenoms()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *HomeController) Statistics(c *gin.Context) {
	resp, err := homeService.Statistics()
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
