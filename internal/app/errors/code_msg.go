package errors

const (
	ErrInvalidParams         = 40000 // 错误的请求参数
	ErrSystemError           = 50000 // 系统异常
	ErrRelayerNoAccessDetail = 50001 // 不支持relayer进入详情页
	ErrLcdNodeError          = 60000 // lcd节点异常
)
