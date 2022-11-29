package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusOK, response.FailBadRequest(err.Error()))
		return
	}

	// todo
	c.JSON(http.StatusOK, response.Success(nil))
}

func (ctl *IbcTxController) FailureStatistics(c *gin.Context) {
	chain := c.Param("chain")
	if chain == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest("parameter chain is required"))
		return
	}
	// todo
	c.JSON(http.StatusOK, response.Success(nil))
}
