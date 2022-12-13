package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/lcd"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
	"sync"
)

type IAddressService interface {
	BaseInfo(chain, address string) (*vo.BaseInfoResp, errors.Error)
	TokenList(chain, address string) (*vo.AddrTokenListResp, errors.Error)
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

func (svc *AddressService) TokenList(chain, address string) (*vo.AddrTokenListResp, errors.Error) {
	cfg, err := chainCfgRepo.FindOneChainInfo(chain)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.WrapBadRequest(fmt.Errorf("invalid chain %s", chain))
		}

		return nil, errors.Wrap(err)
	}

	denomPriceMap := cache.TokenPriceMap()
	authDenomList, err := authDenomRepo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var stakeDenom string
	for _, val := range authDenomList {
		if val.IsStakingToken && val.Chain == chain {
			stakeDenom = val.Denom
		}
	}

	var (
		balanceToken         []vo.AddrToken
		totalValueUnbonding  = decimal.NewFromFloat(0)
		totalValueBalance    = decimal.NewFromFloat(0)
		totalValueDelegation = decimal.NewFromFloat(0)
		totalAmtUnbonding    = decimal.NewFromFloat(0)
		totalAmtDelegation   = decimal.NewFromFloat(0)
	)
	gw := sync.WaitGroup{}
	gw.Add(3)
	go func() {
		defer gw.Done()
		balances, err := lcd.GetBalances(chain, address, cfg.GrpcRestGateway, cfg.LcdApiPath.BalancesPath)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		balanceToken = make([]vo.AddrToken, 0, len(balances.Balances))
		for _, val := range balances.Balances {
			addrToken := vo.AddrToken{
				Denom:                val.Denom,
				Chain:                chain,
				DenomAvaliableAmount: val.Amount,
			}
			denom, err := denomRepo.FindByDenomChain(val.Denom, chain)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				logrus.Error(err.Error())
				continue
			}
			//denom exist in ibc_denom
			if denom != nil {
				//update denom_type,base_denom,base_denom_chain
				addrToken.DenomType = tokenType(authDenomList, denom.BaseDenom, chain)
				if addrToken.DenomType == entity.TokenTypeAuthed && val.Denom == denom.BaseDenom {
					addrToken.DenomType = entity.TokenTypeGenesis
				}
				addrToken.BaseDenom = denom.BaseDenom
				addrToken.BaseDenomChain = denom.BaseDenomChain

				//update denom_value,total_value,price
				if coin, ok := denomPriceMap[denom.BaseDenom+denom.BaseDenomChain]; ok {
					if coin.Scale > 0 {
						addrToken.Price = coin.Price
						decAmt, _ := decimal.NewFromString(val.Amount)
						baseDenomValue := decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
						addrToken.DenomValue = baseDenomValue.String()
						totalValueBalance = totalValueBalance.Add(baseDenomValue)
					}
				}
			} else {
				addrToken.DenomType = entity.TokenTypeOther
				addrToken.BaseDenom = val.Denom
				addrToken.BaseDenomChain = chain
			}

			balanceToken = append(balanceToken, addrToken)
		}

	}()
	go func() {
		defer gw.Done()
		//delegation, err := lcd.GetDelegation(chain, address, cfg.GrpcRestGateway, "/cosmos/staking/v1beta1/delegations/{address}")
		delegation, err := lcd.GetDelegation(chain, address, cfg.GrpcRestGateway, cfg.LcdApiPath.DelegationPath)
		if err != nil {
			logrus.Error(err.Error())
			return
		}

		for _, val := range delegation.DelegationResponses {
			//update denom_value,total_value,price
			if coin, ok := denomPriceMap[val.Balance.Denom+chain]; ok {
				if coin.Scale > 0 {
					decAmt, _ := decimal.NewFromString(val.Balance.Amount)
					baseDenomValue := decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
					totalValueDelegation = totalValueDelegation.Add(baseDenomValue)
					totalAmtDelegation = totalAmtDelegation.Add(decAmt)
				}
			}

		}
	}()
	go func() {
		defer gw.Done()
		//unbonding, err := lcd.GetUnbonding(chain, address, cfg.GrpcRestGateway, "/cosmos/staking/v1beta1/delegators/{address}/unbonding_delegations")
		unbonding, err := lcd.GetUnbonding(chain, address, cfg.GrpcRestGateway, cfg.LcdApiPath.UnbondingPath)
		if err != nil {
			logrus.Error(err.Error())
			return
		}

		for _, val := range unbonding.UnbondingResponses {
			if len(val.Entries) > 0 {
				denomAmount := val.Entries[0].InitialBalance
				//update denom_value,total_value,price
				if coin, ok := denomPriceMap[stakeDenom+chain]; ok {
					if coin.Scale > 0 {
						decAmt, _ := decimal.NewFromString(denomAmount)
						baseDenomValue := decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
						totalValueUnbonding = totalValueUnbonding.Add(baseDenomValue)
						totalAmtUnbonding = totalAmtUnbonding.Add(decAmt)
					}
				}
			}
		}
	}()
	gw.Wait()

	totalValue := totalValueBalance.Add(totalValueDelegation).Add(totalValueUnbonding)
	otherAmt := totalAmtUnbonding.Add(totalAmtDelegation)
	for i, val := range balanceToken {
		if val.Denom == stakeDenom && val.Chain == chain {
			avaliableAmount, _ := decimal.NewFromString(val.DenomAvaliableAmount)
			val.DenomAmount = otherAmt.Add(avaliableAmount).String()
			balanceToken[i] = val
		}
	}
	return &vo.AddrTokenListResp{
		TotalValue: totalValue.String(),
		Tokens:     balanceToken,
	}, nil
}

func tokenType(baseDenomList entity.AuthDenomList, baseDenom, chain string) entity.TokenType {
	for _, v := range baseDenomList {
		if v.Chain == chain && v.Denom == baseDenom {
			return entity.TokenTypeAuthed
		}
	}
	return entity.TokenTypeOther
}
