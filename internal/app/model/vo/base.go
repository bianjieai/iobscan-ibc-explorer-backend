package vo

import "math"

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
