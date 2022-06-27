package rest

import (
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
)

type ChannelController struct {
}

// List channel page 页面
func (ctl *ChannelController) List(c *gin.Context) {
	var req vo.ChannelListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	res, e := channelService.List(&req)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
