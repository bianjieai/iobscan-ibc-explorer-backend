package constant

const (
	EnvNameConfigFilePath = "CONFIG_FILE_PATH"

	DefaultTimezone   = "UTC"
	DefaultTimeFormat = "2006-01-02 15:04:05"
	TimeFormatMMDD    = "0102"
	DefaultCurrency   = "$"

	Cosmos = "cosmos"
	Iris   = "iris"

	MsgTypeTransfer        = "transfer"
	MsgTypeRecvPacket      = "recv_packet"
	MsgTypeTimeoutPacket   = "timeout_packet"
	MsgTypeAcknowledgement = "acknowledge_packet"

	AccountsDailyStatisticName = "accounts_daily"
	TxALlStatisticName         = "tx_all"
	TxFailedStatisticName      = "tx_failed"
)
