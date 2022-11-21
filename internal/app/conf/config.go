package conf

import (
	"bytes"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/enum"
	"github.com/spf13/viper"
)

type Config struct {
	App           App
	Mongo         Mongo
	HintIndexName HintIndexName `mapstructure:"hint_index_name"`
	Mysql         Mysql
	Redis         Redis
	Log           Log
	Spi           Spi
	Task          Task
	ChainConfig   ChainConfig `mapstructure:"chain_config"`
}

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
	TimeZone string `mapstructure:"time_zone"`
}

type Mongo struct {
	Url      string
	Database string
}

type HintIndexName struct {
	MsgsMsgSignerMsgsTypeTimeIndexName string `mapstructure:"msgs_msg_signer_msgs_type_time_index_name"`
	TimeMsgsTypeIndexName              string `mapstructure:"time_msgs_type_index_name"`
	CreateAtIndexName                  string `mapstructure:"create_at_index_name"`
}

type App struct {
	Name                 string
	Addr                 string
	Env                  string
	StartMonitor         bool  `mapstructure:"start_monitor"`
	StartTask            bool  `mapstructure:"start_task"`
	StartOneOffTask      bool  `mapstructure:"start_one_off_task"`
	ApiCacheAliveSeconds int   `mapstructure:"api_cache_alive_seconds"`
	MaxPageSize          int64 `mapstructure:"max_page_size"`
	Version              string
	Prometheus           string `mapstructure:"prometheus_port"`
}

type Redis struct {
	Addrs    string
	User     string
	Password string `json:"-"`
	Mode     enum.RedisMode
	Db       int
}

type Log struct {
	LogLevel           string `mapstructure:"log_level"`
	LogFileName        string `mapstructure:"log_file_name"`
	LogPath            string `mapstructure:"log_path"`
	LogMaxAgeDay       int    `mapstructure:"log_max_age_day"`
	LogRotationTimeDay int    `mapstructure:"log_rotation_time_day"`
	LogOutput          string `mapstructure:"log_output"`
}

type Task struct {
	CronTimeChainTask                 int   `mapstructure:"cron_time_chain_task"`
	CronTimeChannelTask               int   `mapstructure:"cron_time_channel_task"`
	CronTimeRelayerTask               int   `mapstructure:"cron_time_relayer_task"`
	CronTimeStatisticTask             int   `mapstructure:"cron_time_statistic_task"`
	CronTimeTokenTask                 int   `mapstructure:"cron_time_token_task"`
	CronTimeTokenPriceTask            int   `mapstructure:"cron_time_token_price_task"`
	CronTimeChainConfigTask           int   `mapstructure:"cron_time_chain_config_task"`
	CronTimeDenomCalculateTask        int   `mapstructure:"cron_time_denom_calculate_task"`
	CronTimeDenomUpdateTask           int   `mapstructure:"cron_time_denom_update_task"`
	CronTimeSyncTransferTxTask        int   `mapstructure:"cron_time_sync_transfer_tx_task"`
	CronTimeIbcTxRelateTask           int   `mapstructure:"cron_time_ibc_tx_relate_task"`
	CronTimeIbcTxMigrateTask          int   `mapstructure:"cron_time_ibc_tx_migrate_task"`
	RedisLockExpireTime               int   `mapstructure:"redis_lock_expire_time"`
	SingleChainSyncTransferTxMax      int   `mapstructure:"single_chain_sync_transfer_tx_max"`
	SingleChainIbcTxRelateMax         int   `mapstructure:"single_chain_ibc_tx_relate_max"`
	FixDenomTraceDataStartTime        int64 `mapstructure:"fix_denom_trace_data_start_time"`
	FixDenomTraceDataEndTime          int64 `mapstructure:"fix_denom_trace_data_end_time"`
	FixDenomTraceHistoryDataStartTime int64 `mapstructure:"fix_denom_trace_history_data_start_time"`
	FixDenomTraceHistoryDataEndTime   int64 `mapstructure:"fix_denom_trace_history_data_end_time"`
	CronTimeSyncAckTxTask             int   `mapstructure:"cron_time_sync_ack_tx_task"`

	SwitchFixDenomTraceHistoryDataTask bool `mapstructure:"switch_fix_denom_trace_history_data_task"`
	SwitchFixDenomTraceDataTask        bool `mapstructure:"switch_fix_denom_trace_data_task"`
	SwitchAddChainTask                 bool `mapstructure:"switch_add_chain_task"`
	SwitchOnlyInitRelayerData          bool `mapstructure:"switch_only_init_relayer_data"`
	SwitchIbcTxMigrateTask             bool `mapstructure:"switch_ibc_tx_migrate_task"`
	SwitchIbcTokenStatisticsTask       bool `mapstructure:"switch_ibc_token_statistics_task"`
	SwitchIbcChannelStatisticsTask     bool `mapstructure:"switch_ibc_channel_statistics_task"`
	SwitchIbcRelayerStatisticsTask     bool `mapstructure:"switch_ibc_relayer_statistics_task"`
	SwitchAddTransferDataTask          bool `mapstructure:"switch_add_transfer_data_task"`
	SwitchFixFailRecvPacketTask        bool `mapstructure:"switch_fix_fail_recv_packet_task"`
	SwitchFixDcChainIdTask             bool `mapstructure:"switch_fix_dc_chain_id_task"`
	SwitchFixBaseDenomChainIdTask      bool `mapstructure:"switch_fix_base_denom_chain_id_task"`

	SyncTransferTxWorkerNum    int `mapstructure:"sync_transfer_tx_worker_num"`
	IbcTxRelateWorkerNum       int `mapstructure:"ibc_tx_relate_worker_num"`
	FixDenomTraceDataWorkerNum int `mapstructure:"fix_denom_trace_data_worker_num"`

	CreateAtUseTxTime bool `mapstructure:"create_at_use_tx_time"`
}

type Spi struct {
	CoingeckoPriceUrl string `mapstructure:"coingecko_price_url"`
}

type ChainConfig struct {
	NewChains         string `mapstructure:"new_chains"`
	AddTransferChains string `mapstructure:"add_transfer_chains"`
}

func ReadConfig(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
