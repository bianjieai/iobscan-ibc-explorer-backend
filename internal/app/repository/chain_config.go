package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ChainConfigFieldCurrentChainId  = "current_chain_id"
	ChainConfigFieldGrpcRestGateway = "grpc_rest_gateway"
	ChainConfigFieldChainName       = "chain_name"
	ChainConfigFieldPrettyName      = "pretty_name"
	ChainConfigFieldAddrPrefix      = "addr_prefix"
	ChainConfigFieldIcon            = "icon"
	ChainConfigFieldStatus          = "status"
	ChainConfigFieldIbcInfo         = "ibc_info"
	ChainConfigFieldIbcInfoHashLcd  = "ibc_info_hash_lcd"
	ChainConfigFieldLcdApiPath      = "lcd_api_path"
)

type IChainConfigRepo interface {
	FindAll() ([]*entity.ChainConfig, error)
	FindAllChainInfos() ([]*entity.ChainConfig, error)
	FindAllOpenChainInfos() ([]*entity.ChainConfig, error)
	FindOne(chain string) (*entity.ChainConfig, error)
	FindOneChainInfo(chain string) (*entity.ChainConfig, error)
	UpdateIbcInfo(config *entity.ChainConfig) error
	UpdateLcdApi(config *entity.ChainConfig) error
	Count() (int64, error)
}

var _ IChainConfigRepo = new(ChainConfigRepo)

type ChainConfigRepo struct {
}

func (repo *ChainConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ChainConfig{}.CollectionName())
}

func (repo *ChainConfigRepo) Count() (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{}).Count()
}

func (repo *ChainConfigRepo) FindAll() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindAllChainInfos() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{}).
		Select(bson.M{ChainConfigFieldCurrentChainId: 1, ChainConfigFieldChainName: 1, ChainConfigFieldPrettyName: 1, ChainConfigFieldIcon: 1, ChainConfigFieldGrpcRestGateway: 1,
			ChainConfigFieldLcdApiPath: 1, ChainConfigFieldStatus: 1, ChainConfigFieldAddrPrefix: 1}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindAllOpenChainInfos() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{ChainConfigFieldStatus: entity.ChainStatusOpen}).
		Select(bson.M{ChainConfigFieldCurrentChainId: 1, ChainConfigFieldChainName: 1, ChainConfigFieldPrettyName: 1, ChainConfigFieldIcon: 1, ChainConfigFieldGrpcRestGateway: 1,
			ChainConfigFieldLcdApiPath: 1}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindOne(chain string) (*entity.ChainConfig, error) {
	var res *entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{ChainConfigFieldChainName: chain}).One(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindOneChainInfo(chain string) (*entity.ChainConfig, error) {
	var res *entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{ChainConfigFieldChainName: chain}).
		Select(bson.M{ChainConfigFieldCurrentChainId: 1, ChainConfigFieldChainName: 1, ChainConfigFieldPrettyName: 1, ChainConfigFieldIcon: 1, ChainConfigFieldGrpcRestGateway: 1,
			ChainConfigFieldLcdApiPath: 1}).
		One(&res)
	return res, err
}

func (repo *ChainConfigRepo) UpdateIbcInfo(config *entity.ChainConfig) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainConfigFieldChainName: config.ChainName}, bson.M{
		"$set": bson.M{
			ChainConfigFieldIbcInfo:        config.IbcInfo,
			ChainConfigFieldIbcInfoHashLcd: config.IbcInfoHashLcd,
		}})
}

func (repo *ChainConfigRepo) UpdateLcdApi(config *entity.ChainConfig) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainConfigFieldChainName: config.ChainName}, bson.M{
		"$set": bson.M{
			ChainConfigFieldGrpcRestGateway:  config.GrpcRestGateway,
			"lcd_api_path.channels_path":     config.LcdApiPath.ChannelsPath,
			"lcd_api_path.client_state_path": config.LcdApiPath.ClientStatePath,
			//"lcd_api_path.supply_path":       config.LcdApiPath.SupplyPath,
			//"lcd_api_path.balances_path":     config.LcdApiPath.BalancesPath,
			//"lcd_api_path.params_path":       config.LcdApiPath.ParamsPath,
		}})
}
