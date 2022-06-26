package vo

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageInfo struct {
	TotalItem int   `json:"total_item"`
	TotalPage int   `json:"total_page"`
	PageNum   int64 `json:"page_num"`
	PageSize  int64 `json:"page_size"`
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
	return (pageNum - 1) * pageSize, pageSize
}
