package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainConfigRepo interface {
	FindAll() ([]*entity.ChainConfig, error)
	FindAllChainInfos() ([]*entity.ChainConfig, error)
	FindAllOpenChainInfos() ([]*entity.ChainConfig, error)
	FindOne(chainId string) (*entity.ChainConfig, error)
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
		Select(bson.M{"chain_id": 1, "chain_name": 1, "pretty_name": 1, "icon": 1, "lcd": 1, "lcd_api_path": 1, "status": 1}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindAllOpenChainInfos() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{"status": entity.ChainStatusOpen}).
		Select(bson.M{"chain_id": 1, "chain_name": 1, "pretty_name": 1, "icon": 1, "lcd": 1, "lcd_api_path": 1}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindOne(chainId string) (*entity.ChainConfig, error) {
	var res *entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{"chain_id": chainId}).One(&res)
	return res, err
}

func (repo *ChainConfigRepo) UpdateIbcInfo(config *entity.ChainConfig) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"chain_id": config.ChainId}, bson.M{
		"$set": bson.M{
			"ibc_info":          config.IbcInfo,
			"ibc_info_hash_lcd": config.IbcInfoHashLcd,
		}})
}

func (repo *ChainConfigRepo) UpdateLcdApi(config *entity.ChainConfig) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"chain_id": config.ChainId}, bson.M{
		"$set": bson.M{
			"lcd":                            config.Lcd,
			"lcd_api_path.channels_path":     config.LcdApiPath.ChannelsPath,
			"lcd_api_path.client_state_path": config.LcdApiPath.ClientStatePath,
			//"lcd_api_path.supply_path":       config.LcdApiPath.SupplyPath,
			//"lcd_api_path.balances_path":     config.LcdApiPath.BalancesPath,
			//"lcd_api_path.params_path":       config.LcdApiPath.ParamsPath,
		}})
}
