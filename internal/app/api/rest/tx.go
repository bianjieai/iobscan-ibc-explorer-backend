package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
)

type IbcTxController struct {
}

func (ctl *IbcTxController) Query(c *gin.Context) {
	txHash := c.Param("tx_hash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter tx_hash is required"))
		return
	}

	var req vo.TxReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailBadRequest(err.Error()))
		return
	}

	if req.Chain == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter chain is required"))
		return
	}

	res, e := txService.Query(txHash, req)
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(res))
}

func (ctl *IbcTxController) FailureStatistics(c *gin.Context) {
	chain := c.Param("chain")
	if chain == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter chain is required"))
		return
	}

	var req vo.FailureStatisticsReq
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
	}

	if startTime >= endTime {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("end_date must be greater than start_date"))
		return
	}

	res, e := txService.FailureStatistics(chain, startTime, endTime)
	if e != nil {
		c.JSON(response.HttpCode(e), response.FailError(e))
		return
	}

	c.JSON(http.StatusOK, response.Success(res))
}
