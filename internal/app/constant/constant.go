package constant

const (
	EnvNameConfigFilePath = "CONFIG_FILE_PATH"

	DefaultTimezone       = "UTC"
	DefaultTimeFormat     = "2006-01-02 15:04:05"
	TimeFormatMMDD        = "0102"
	DefaultCurrency       = "$"
	UnknownTokenPrice     = -1
	UnknownDenomAmount    = "-1"
	ZeroDenomAmount       = "0"
	IBCTokenPreFix        = "ibc"
	IBCHopsIndex          = "/channel"
	DefaultValuePrecision = 5
	ChannelStateOpen      = "STATE_OPEN"
	DefaultPageSize       = 10
	DefaultPageNum        = 1
	OtherDenom            = "others"
	AllChain              = "allchain"
	Cosmos                = "cosmos"
	Iris                  = "iris"

	MsgTypeRecvPacket      = "recv_packet"
	MsgTypeTimeoutPacket   = "timeout_packet"
	MsgTypeAcknowledgement = "acknowledge_packet"
	MsgTypeUpdateClient    = "update_client"

	ChannelOpenStatisticName  = "channel_opened"
	ChannelCloseStatisticName = "channel_closed"
	ChannelAllStatisticName   = "channel_all"
	Channel24hStatisticName   = "channels_24hr"
	TxALlStatisticName        = "tx_all"
	TxFailedStatisticName     = "tx_failed"
)
