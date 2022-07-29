package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
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
