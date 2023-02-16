package rest

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type IbcFeeController struct {
}

func (ctl *IbcFeeController) ChinFeeStatistics(c *gin.Context) {
	chain := c.Param("chain")
	if chain == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter chain is required"))
		return
	}

	var req vo.ChainFeeStatisticsReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailBadRequest(err.Error()))
		return
	}

	var startTime, endTime int64
	if req.StartDate != "" {
		startTimeStr := fmt.Sprintf("%s %s", req.StartDate, "00:00:00")
		startTimeParse, err := time.Parse(constant.DefaultTimeFormat, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailBadRequest("invalid start_date"))
			return
		}
		startTime = startTimeParse.Unix()
	}

	if req.EndDate != "" {
		endTimeStr := fmt.Sprintf("%s %s", req.EndDate, "23:59:59")
		endTimeParse, err := time.Parse(constant.DefaultTimeFormat, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailBadRequest("invalid end_date"))
			return
		}
		endTime = endTimeParse.Unix()
	} else {
		_, endTime = utils.YesterdayUnix()
	}

	if startTime > endTime {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("end_date must be greater than start_date"))
		return
	}

	res, e := feeService.ChainFeeStatistics(chain, startTime, endTime)
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(res))
}
