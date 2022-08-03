package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
)

type ChainController struct {
}

// @Summary list
// @Description get IBC BUSD of chains
// @ID list
// @Tags app_version
// @Accept  json
// @Produce  json
// @Param   page_num    query   int true    "page num" Default(1)
// @Param   page_size   query   int true    "page size" Default(10)
// @Param   use_count   query   bool false    "if used count" Enums(true, false)
// @Success 200 {object} vo.ChainListResp	"success"
// @Router /ibc/chainList/api_support [get]
func (ctl *ChainController) List(c *gin.Context) {
	var req vo.ChainListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailError(errors.Wrap(err)))
		return
	}
	if req.UseCount {
		total, err := chainService.Count()
		if err != nil {
			c.JSON(http.StatusOK, response.FailError(err))
			return
		}
		c.JSON(http.StatusOK, response.Success(total))
		return
	}
	resp, err := chainService.List(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
