package integration

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/redis"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	chainService                  service.IChainService
	txRepo                        repository.ITxRepo
	chainConfigRepo               repository.IChainConfigRepo
	chainRegistryRepo             repository.IChainRegistryRepo
	ibcTxFailLogRepo              repository.IIBCTxFailLogRepo
	ibcChainInflowStatisticsRepo  repository.IChainInflowStatisticsRepo
	ibcChainOutflowStatisticsRepo repository.IChainOutflowStatisticsRepo
	ibcChainFeeStatisticsRepo     repository.IChainFeeStatisticsRepo
	ibcRelayerFeeStatisticsRepo   repository.IRelayerFeeStatisticsRepo
	suite.Suite
}

type SubTest struct {
	testName string
	testCase func(s IntegrationTestSuite)
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IntegrationTestSuite))
}

var (
	mgo        *qmgo.Client
	redisConn  *redis.Client
	ctx                            = context.Background()
	txService  service.ITxService  = new(service.TxService)
	feeService service.IFeeService = new(service.FeeService)
)

func (s *IntegrationTestSuite) SetupSuite() {
	cfg := testConfig()
	global.Config = cfg

	var maxPoolSize uint64 = 4096
	var timeout = int64(600 * time.Second)
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri: cfg.Mongo.Url,
		ReadPreference: &qmgo.ReadPref{
			MaxStalenessMS: 90000,
			Mode:           readpref.SecondaryPreferredMode,
		},
		Database:        cfg.Mongo.Database,
		MaxPoolSize:     &maxPoolSize,
		SocketTimeoutMS: &timeout,
	})
	if err != nil {
		logrus.Fatalf("connect mongo failed, uri: %s, err:%s", cfg.Mongo.Url, err.Error())
	}
	mgo = client
	redisConn = cache.InitRedisClient(cfg.Redis)
	repository.InitMgo(cfg.Mongo, ctx)

	s.chainService = new(service.ChainService)
	s.txRepo = new(repository.TxRepo)
	s.chainConfigRepo = new(repository.ChainConfigRepo)
	s.chainRegistryRepo = new(repository.ChainRegistryRepo)
	s.ibcTxFailLogRepo = new(repository.IBCTxFailLogRepo)
	s.ibcChainInflowStatisticsRepo = new(repository.ChainInflowStatisticsRepo)
	s.ibcChainOutflowStatisticsRepo = new(repository.ChainOutflowStatisticsRepo)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	mgo.Close(ctx)
	redisConn.Close()
	s.T().Log("close all db connect")
}

func testConfig() *conf.Config {
	data, err := ioutil.ReadFile("../configs/cfg.toml")
	if err != nil {
		panic(err)
	}
	cfg, err := conf.ReadConfig(data)
	if err != nil {
		panic(err)
	}
	return cfg
}
