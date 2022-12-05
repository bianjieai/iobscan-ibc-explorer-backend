package response

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"net/http"
)

var ErrHttpCodeMap = map[int]int{
	errors.ErrBadRequest:  http.StatusBadRequest,
	errors.ErrTxNotFound:  http.StatusNotFound,
	errors.ErrTxNotUnique: http.StatusBadRequest,
	errors.ErrSystemError: http.StatusInternalServerError,
}

func HttpCode(err errors.Error) int {
	code, ok := ErrHttpCodeMap[err.Code()]
	if !ok {
		return http.StatusInternalServerError
	}
	return code
}

func Success(data interface{}) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func FailSystemError() vo.BaseResponse {
	return vo.BaseResponse{
		Code:    errors.ErrSystemError,
		Message: "System error",
	}
}

func FailBadRequest(msg string) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    errors.ErrBadRequest,
		Message: msg,
	}
}

func FailError(err errors.Error) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    err.Code(),
		Message: err.Msg(),
	}
}
