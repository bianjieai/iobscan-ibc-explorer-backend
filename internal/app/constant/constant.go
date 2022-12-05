package constant

const (
	EnvNameConfigFilePath = "CONFIG_FILE_PATH"

	DefaultTimezone       = "UTC"
	DefaultTimeFormat     = "2006-01-02 15:04:05"
	DateFormat            = "2006-01-02"
	TimeFormatMMDD        = "0102"
	DefaultCurrency       = "$"
	UnknownTokenPrice     = -1
	UnknownDenomAmount    = "-1"
	ZeroDenomAmount       = "0"
	IBCTokenPrefix        = "ibc"
	IBCHopsIndex          = "/channel"
	DefaultValuePrecision = 5
	ChannelStateOpen      = "STATE_OPEN"
	DefaultPageSize       = 10
	DefaultPageNum        = 1
	OtherDenom            = "others"
	AllChain              = "allchain"
	Cosmos                = "cosmos"
	Iris                  = "iris"
	PortTransfer          = "transfer"
	DefaultUnboundTime    = 1209600

	DefaultLimit = 500
	IncreHeight  = 5000

	DisplayIbcRecordMax = 500000

	MsgTypeTransfer           = "transfer"
	MsgTypeRecvPacket         = "recv_packet"
	MsgTypeTimeoutPacket      = "timeout_packet"
	MsgTypeAcknowledgement    = "acknowledge_packet"
	MsgTypeUpdateClient       = "update_client"
	MsgTypeChannelOpenConfirm = "channel_open_confirm"

	ChannelOpenStatisticName  = "channel_opened"
	ChannelCloseStatisticName = "channel_closed"
	ChannelAllStatisticName   = "channel_all"
	Channel24hStatisticName   = "channels_24hr"
	Chains24hStatisticName    = "chains_24hr"
	ChainsAllStatisticName    = "chain_all"
	Tx24hAllStatisticName     = "tx_24hr_all"
	TxLatestAllStatisticName  = "tx_latest_all"
	TxAllStatisticName        = "tx_all"
	TxSuccessStatisticName    = "tx_success"
	TxFailedStatisticName     = "tx_failed"
	BaseDenomAllStatisticName = "base_denom_all"
	DenomAllStatisticName     = "denom_all"

	IbcCoreConnectionUri = "%s/ibc/core/connection/%s/client_connections/%s"
	IbcCoreChannelsUri   = "%s/ibc/core/channel/%s/connections/%s/channels"

	HeaderTimestamp = "X-Timestamp"
	HeaderSignature = "X-Signature"
	HeaderApiKey    = "X-Api-Key"
	NetworkDelay    = 15

	ContactUs = "If you don't find the chain you need, you can contact us by twitter(https://twitter.com/iobscan_ibc)"
)

var HomeStatistics = []string{
	ChannelOpenStatisticName, ChannelCloseStatisticName, ChannelAllStatisticName, Channel24hStatisticName,
	Chains24hStatisticName, ChainsAllStatisticName,
	Tx24hAllStatisticName, TxAllStatisticName, TxSuccessStatisticName, TxFailedStatisticName,
	BaseDenomAllStatisticName, DenomAllStatisticName,
}

var RelayerDetailTxsType = []string{MsgTypeRecvPacket, MsgTypeAcknowledgement, MsgTypeTimeoutPacket}

const (
	UnAuth = "Others"
	//AllChain = "allchain"
)
