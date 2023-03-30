package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type TaskController struct {
}

func (ctl *TaskController) Run(c *gin.Context) {
	taskName := c.Param("task_name")
	if taskName == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter task_name is required"))
		return
	}
	lockKey := fmt.Sprintf("%s:%s", "TaskController", taskName)
	if err := cache.GetRedisClient().Lock(lockKey, time.Now().Unix(), time.Hour); err != nil {
		c.JSON(http.StatusTooManyRequests, response.FailBadRequest("Please try again later"))
		return
	}

	go func() {
		st := time.Now().Unix()
		res := 0
		logrus.Infof("TaskController task %s start", taskName)

		switch taskName {
		case ibcTxFailLogTask.Name():
			var req vo.TaskReq
			if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
				logrus.Errorf("TaskController run %s err, %v", taskName, err)
				return
			}
			res = ibcTxFailLogTask.RunWithParam(req.StartTime, req.EndTime, req.IsTargetHistory)
		case iBCChainFeeStatisticTask.Name():
			chain := c.PostForm("chain")
			if chain == "" {
				iBCChainFeeStatisticTask.RunAllChain()
			} else {
				startTime, err := strconv.ParseInt(c.PostForm("start_time"), 10, 64)
				if err != nil {
					logrus.Errorf("TaskController run %s err, %v", taskName, err)
					return
				}
				endTime, err := strconv.ParseInt(c.PostForm("end_time"), 10, 64)
				if err != nil {
					logrus.Errorf("TaskController run %s err, %v", taskName, err)
					return
				}
				res = iBCChainFeeStatisticTask.RunWithParam(chain, startTime, endTime)
			}
		case ibcAddressStatisticTask.Name():
			var req vo.TaskAddressStatisticReq
			if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
				logrus.Errorf("TaskController run %s err, %v", taskName, err)
				return
			}
			if req.Chain == "" {
				iBCChainFeeStatisticTask.RunAllChain()
			} else {
				res = ibcAddressStatisticTask.RunWithParam(req.Chain, req.StartTime, req.EndTime)
			}
		default:
			logrus.Errorf("TaskController run %s err, %s", taskName, "unknown task")
		}

		logrus.Infof("TaskController task %s end, time use %d(s), exec status: %d", taskName, time.Now().Unix()-st, res)
	}()
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, response.Success("task is running"))

}
