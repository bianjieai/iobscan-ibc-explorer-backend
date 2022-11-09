package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"time"
)

var (
	tokenRepo                  repository.ITokenRepo                  = new(repository.TokenRepo)
	tokenStatisticsRepo        repository.ITokenTraceRepo             = new(repository.TokenTraceRepo)
	channelRepo                repository.IChannelRepo                = new(repository.ChannelRepo)
	denomRepo                  repository.IDenomRepo                  = new(repository.DenomRepo)
	chainRepo                  repository.IChainRepo                  = new(repository.IbcChainRepo)
	relayerRepo                repository.IRelayerRepo                = new(repository.IbcRelayerRepo)
	statisticRepo              repository.IStatisticRepo              = new(repository.IbcStatisticRepo)
	chainCfgRepo               repository.IChainConfigRepo            = new(repository.ChainConfigRepo)
	ibcTxRepo                  repository.IExIbcTxRepo                = new(repository.ExIbcTxRepo)
	txRepo                     repository.ITxRepo                     = new(repository.TxRepo)
	exSearchRecordRepo         repository.IExSearchRecordRepo         = new(repository.ExSearchRecordRepo)
	relayerAddrChannelRepo     repository.IRelayerAddressChannelRepo  = new(repository.RelayerAddressChannelRepo)
	relayerDenomStatisticsRepo repository.IRelayerDenomStatisticsRepo = new(repository.RelayerDenomStatisticsRepo)
	relayerTxsCache            cache.RelayerTotalTxsCacheRepo
	relayerTrendCache          cache.RelayerRelayedTrendCacheRepo
	lcdTxDataCache             cache.LcdTxDataCacheRepo
	lcdAddrCache               cache.LcdAddrCacheRepo
	relayerCfgRepo             repository.IRelayerConfigRepo = new(cache.RelayerConfigCacheRepo)
	baseDenomRepo              cache.BaseDenomCacheRepo
)

type (
	LcdTxData struct {
		TxResponse struct {
			Logs []LogData `json:"logs"`
			Tx   struct {
				Body struct {
					Messages []LcdMessage `json:"messages"`
				} `json:"body"`
			} `json:"tx"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"tx_response"`
	}
	LogData struct {
		MsgIndex int            `json:"msg_index"`
		Log      string         `json:"log"`
		Events   []entity.Event `json:"events"`
	}
	LcdMessage struct {
		Type            string      `json:"@type"`
		Packet          interface{} `json:"packet,omitempty"`
		ProofCommitment string      `json:"proof_commitment,omitempty"`
		ProofHeight     interface{} `json:"proof_height,omitempty"`
		Signer          string      `json:"signer,omitempty"`

		SourcePort       string      `json:"source_port,omitempty"`
		SourceChannel    string      `json:"source_channel,omitempty"`
		Token            interface{} `json:"token,omitempty"`
		Sender           string      `json:"sender,omitempty"`
		Receiver         string      `json:"receiver,omitempty"`
		TimeoutHeight    interface{} `json:"timeout_height,omitempty"`
		TimeoutTimestamp string      `json:"timeout_timestamp,omitempty"`

		ProofUnreceived  string `json:"proof_unreceived,omitempty"`
		NextSequenceRecv string `json:"next_sequence_recv,omitempty"`
		Acknowledgement  string `json:"acknowledgement,omitempty"`
		ProofAcked       string `json:"proof_acked,omitempty"`
	}
	LcdErrRespond struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Details []interface{} `json:"details"`
	}
)
