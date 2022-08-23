package rest

import (
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
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

	var res interface{}
	var err errors.Error
	if req.UseCount {
		res, err = tokenService.ListCount(&req)
	} else {
		res, err = tokenService.List(&req)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
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

	var res interface{}
	var err errors.Error
	if req.UseCount {
		res, err = tokenService.IBCTokenListCount(baseDenom, &req)
	} else {
		res, err = tokenService.IBCTokenList(baseDenom, &req)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
