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

func SuccessWithMsg(msg string, data interface{}) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    0,
		Message: msg,
		Data:    data,
	}
}

func Fail(code int, msg string, data interface{}) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func FailMsg(msg string) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    errors.ErrSystemError,
		Message: msg,
		Data:    nil,
	}
}

func FailError(err errors.Error) vo.BaseResponse {
	return vo.BaseResponse{
		Code:    err.Code(),
		Message: err.Msg(),
	}
}

func FailBadRequest(err error) vo.BaseResponse {
	return FailError(errors.WrapBadRequest(err))
}
