package vo

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageInfo struct {
	TotalItem int `json:"total_item"`
	TotalPage int `json:"total_page"`
	PageNum   int `json:"page_num"`
	PageSize  int `json:"page_size"`
}
