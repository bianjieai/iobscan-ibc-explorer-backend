package rest

import (
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
}

// List token page 页面
func (ctl *TokenController) List(c *gin.Context) {
	var req vo.TokenListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	res, e := tokenService.List(&req)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}

// IBCTokenList token page 子页面
func (ctl *TokenController) IBCTokenList(c *gin.Context) {
	baseDenom := c.Param("base_denom")
	var req vo.IBCTokenListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	res, e := tokenService.IBCTokenList(baseDenom, &req)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
