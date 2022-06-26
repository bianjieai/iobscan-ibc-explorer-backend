package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
)

type RelayerController struct {
}

func (ctl *RelayerController) List(c *gin.Context) {
	var req vo.RelayerListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	resp, err := relayerService.List(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
