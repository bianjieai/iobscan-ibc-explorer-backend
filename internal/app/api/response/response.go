package response

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

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
