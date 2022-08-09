package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"math"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageInfo struct {
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
	PageNum   int64 `json:"page_num"`
	PageSize  int64 `json:"page_size"`
}

func BuildPageInfo(totalItem, pageNum, pageSize int64) PageInfo {
	p := PageInfo{
		TotalItem: totalItem,
		TotalPage: 0,
		PageNum:   pageNum,
		PageSize:  pageSize,
	}

	if totalItem == 0 {
		return p
	}

	p.TotalPage = int64(math.Ceil(float64(totalItem) / float64(pageSize)))
	return p
}

type Page struct {
	PageNum  int64 `json:"page_num" form:"page_num" binding:"required"`
	PageSize int64 `json:"page_size" form:"page_size" binding:"required"`
}

// parse page param
// get skip, limit variable which used in database
func ParseParamPage(pageNum int64, pageSize int64) (skip int64, limit int64) {
	if pageNum == 0 && pageSize == 0 {
		pageSize = 10
	}
	//limit max pagesize
	if pageSize > constant.MaxPageSize {
		pageSize = constant.MaxPageSize
	}
	return (pageNum - 1) * pageSize, pageSize
}
