package service

import (
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/lcd"
	"github.com/qiniu/qmgo"
)

type IAddressService interface {
	BaseInfo(chain, address string) (*vo.BaseInfoResp, errors.Error)
}

type AddressService struct {
}

var _ IAddressService = new(AddressService)

func PubKeyAlgorithm(pubKeyType string) string {
	if strings.Contains(pubKeyType, constant.ETH_SECP256K1) {
		return constant.ETH_SECP256K1
	}

	if strings.Contains(pubKeyType, constant.SECP256K1) {
		return constant.SECP256K1
	}

	split := strings.Split(pubKeyType, ".")
	if len(split) >= 3 {
		return split[2]
	}

	return ""
}

func (svc *AddressService) BaseInfo(chain, address string) (*vo.BaseInfoResp, errors.Error) {
	cfg, err := chainCfgRepo.FindOneChainInfo(chain)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.WrapBadRequest(fmt.Errorf("invalid chain %s", chain))
		}

		return nil, errors.Wrap(err)
	}

	account, err := lcd.GetAccount(chain, address, cfg.GrpcRestGateway, cfg.LcdApiPath.AccountsPath)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &vo.BaseInfoResp{
		Address:         address,
		PubKey:          account.Account.PubKey.Key,
		Chain:           chain,
		AccountNumber:   account.Account.AccountNumber,
		PubKeyType:      account.Account.PubKey.Type,
		PubKeyAlgorithm: PubKeyAlgorithm(account.Account.PubKey.Type),
		AccountSequence: account.Account.Sequence,
	}, nil
}
