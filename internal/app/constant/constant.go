package constant

const (
	EnvNameZkServices   = "ZK_SERVICES"
	EnvNameZkUsername   = "ZK_USERNAME"
	EnvNameZkPasswd     = "ZK_PASSWD"
	EnvNameZkConfigPath = "ZK_CONFIG_PATH"

	DefaultTimezone       = "UTC"
	DefaultTimeFormat     = "2006-01-02 15:04:05"
	DefaultCurrency       = "$"
	UnknownTokenPrice     = -1
	UnknownDenomAmount    = ""
	ZeroDenomAmount       = "0"
	IBCTokenPreFix        = "ibc"
	IBCHopsIndex          = "/channel"
	DefaultValuePrecision = 5
	ChannelStateOpen      = "STATE_OPEN"
)
