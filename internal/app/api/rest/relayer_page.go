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
