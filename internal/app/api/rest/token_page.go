package rest

import (
	"net/http"
	"strconv"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
}

// List token page 页面
func (ctl *TokenController) List(c *gin.Context) {
	baseDenom := c.Query("base_denom")
	chain := c.Query("chain")
	TokenType := c.Query("token_type")
	userCountStr := c.Query("use_count")
	pageNumStr := c.Query("page_num")
	pageSizeStr := c.Query("page_size")

	useCount, err := strconv.ParseBool(userCountStr)
	if err != nil {
		useCount = false
	}

	pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = constant.DefaultPageNum
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = constant.DefaultPageSize
	}

	res, e := tokenService.List(baseDenom, chain, entity.TokenType(TokenType), useCount, pageNum, pageSize)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}

// IBCTokenList token page 子页面
func (ctl *TokenController) IBCTokenList(c *gin.Context) {
	baseDenom := c.Param("base_denom")
	chain := c.Query("chain")
	TokenType := c.Query("token_type")
	userCountStr := c.Query("use_count")
	pageNumStr := c.Query("page_num")
	pageSizeStr := c.Query("page_size")

	useCount, err := strconv.ParseBool(userCountStr)
	if err != nil {
		useCount = false
	}

	pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = constant.DefaultPageNum
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = constant.DefaultPageSize
	}

	res, e := tokenService.IBCTokenList(baseDenom, chain, entity.TokenStatisticsType(TokenType), useCount, pageNum, pageSize)
	if e != nil {
		c.JSON(http.StatusOK, response.FailError(e))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}
