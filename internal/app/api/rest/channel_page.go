package rest

import (
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ChannelController struct {
}

// List token page 页面
func (ctl *ChannelController) List(c *gin.Context) {
	chainStr := c.Query("chain")
	statusStr := c.Query("status")
	userCountStr := c.Query("use_count")
	pageNumStr := c.Query("page_num")
	pageSizeStr := c.Query("page_size")

	var status int
	if statusStr == "" {
		status = 0
	} else {
		status, _ = strconv.Atoi(statusStr)
	}

	var chainA, chainB string
	if chainStr == "" {
		chainA = constant.AllChain
		chainB = constant.AllChain
	} else {
		split := strings.Split(chainStr, ",")
		if len(split) != 2 {
			c.JSON(http.StatusOK, response.FailBadRequest(fmt.Errorf("chain format error")))
			return
		} else {
			chainA = split[0]
			chainB = split[1]
		}
	}

	useCount, err := strconv.ParseBool(userCountStr)
	if err != nil {
		useCount = false
	}

	pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = constant.DefaultPageNum
	}

	res, e := channelService.List(chainA, chainB, entity.ChannelStatus(status), useCount, pageNum, pageSize)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
