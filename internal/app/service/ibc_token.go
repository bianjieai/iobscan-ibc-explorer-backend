package service

import (
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type ITokenService interface {
	List(baseDenom, chainId string, tokenType entity.TokenType, useCount bool, pageNum, pageSize int64) (*vo.TokenListResp, errors.Error)
	IBCTokenList(baseDenom, chainId string, tokenType entity.TokenStatisticsType, useCount bool, pageNum, pageSize int64) (*vo.IBCTokenListResp, errors.Error)
}

type TokenService struct {
}

var _ ITokenService = new(TokenService)

func (svc *TokenService) List(baseDenom, chainId string, tokenType entity.TokenType, useCount bool, pageNum, pageSize int64) (*vo.TokenListResp, errors.Error) {
	var baseDenomList []string
	if baseDenom != "" {
		if strings.ToLower(baseDenom) == constant.OtherDenom {
			others, err := denomRepo.FindTokenOthers()
			if err != nil {
				return nil, errors.Wrap(err)
			}

			for _, v := range others {
				baseDenomList = append(baseDenomList, v.BaseDenom)
			}
		} else {
			baseDenomList = []string{baseDenom}
		}
	}

	list, err := tokenRepo.List(baseDenomList, chainId, tokenType, pageNum, pageSize)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.TokenItem, 0, len(list))
	for _, v := range list {
		items = append(items, vo.TokenItem{
			BaseDenom:         v.BaseDenom,
			ChainId:           v.ChainId,
			Supply:            v.Supply,
			Currency:          v.Currency,
			Price:             v.Price,
			ChainsInvolved:    v.ChainsInvolved,
			IBCTransferTxs:    v.TransferTxs,
			IBCTransferAmount: v.TransferAmount,
		})
	}

	var totalItem int64
	if useCount {
		totalItem, err = tokenRepo.CountList(baseDenomList, chainId, tokenType)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, pageNum, pageSize)
	return &vo.TokenListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}

func (svc *TokenService) IBCTokenList(baseDenom, chainId string, tokenType entity.TokenStatisticsType, useCount bool, pageNum, pageSize int64) (*vo.IBCTokenListResp, errors.Error) {
	list, err := tokenStatisticsRepo.List(baseDenom, chainId, tokenType, pageNum, pageSize)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.IBCTokenItem, 0, len(list))
	for _, v := range list {
		items = append(items, vo.IBCTokenItem{
			Denom:      v.Denom,
			DenomPath:  v.DenomPath,
			ChainId:    v.ChainId,
			TokenType:  v.Type,
			IBCHops:    v.IBCHops,
			Amount:     v.DenomAmount,
			ReceiveTxs: v.ReceiveTxs,
		})
	}

	var totalItem int64
	if useCount {
		totalItem, err = tokenStatisticsRepo.CountList(baseDenom, chainId, tokenType)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, pageNum, pageSize)
	return &vo.IBCTokenListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}
