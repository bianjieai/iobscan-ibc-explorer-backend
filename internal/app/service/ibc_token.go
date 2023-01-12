package service

import (
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

type ITokenService interface {
	List(req *vo.TokenListReq) (*vo.TokenListResp, errors.Error)
	ListCount(req *vo.TokenListReq) (int64, errors.Error)
	IBCTokenList(req *vo.IBCTokenListReq) (*vo.IBCTokenListResp, errors.Error)
	IBCTokenListCount(req *vo.IBCTokenListReq) (int64, errors.Error)
}

type TokenService struct {
}

var _ ITokenService = new(TokenService)

func (svc *TokenService) List(req *vo.TokenListReq) (*vo.TokenListResp, errors.Error) {
	if req.BaseDenomChain != "" && req.Chain != "" && req.BaseDenomChain != req.Chain {
		return &vo.TokenListResp{
			Items: make([]vo.TokenItem, 0, 0),
		}, nil
	}

	var chain string
	if req.BaseDenomChain != "" {
		chain = req.BaseDenomChain
	} else {
		chain = req.Chain
	}

	baseDenomList, err := svc.analyzeBaseDenom(req.BaseDenom)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	list, err := tokenRepo.List(baseDenomList, chain, req.TokenType, skip, limit)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.TokenItem, 0, len(list))
	for _, v := range list {
		items = append(items, vo.TokenItem{
			BaseDenom:         v.BaseDenom,
			Chain:             v.Chain,
			TokenType:         v.Type,
			Supply:            v.Supply,
			Currency:          v.Currency,
			Price:             v.Price,
			ChainsInvolved:    v.ChainsInvolved,
			IBCTransferTxs:    v.TransferTxs,
			IBCTransferAmount: v.TransferAmount,
		})
	}

	return &vo.TokenListResp{
		Items: items,
	}, nil
}

func (svc *TokenService) ListCount(req *vo.TokenListReq) (int64, errors.Error) {
	if req.BaseDenomChain != "" && req.Chain != "" && req.BaseDenomChain != req.Chain {
		return 0, nil
	}

	baseDenomList, err := svc.analyzeBaseDenom(req.BaseDenom)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	totalItem, err := tokenRepo.CountList(baseDenomList, req.Chain, req.TokenType)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return totalItem, nil
}

func (svc *TokenService) analyzeBaseDenom(baseDenom string) ([]string, error) {
	if baseDenom == "" {
		return nil, nil
	}

	var baseDenomList []string
	if strings.ToLower(baseDenom) == constant.OtherDenom {
		others, err := getUnAuthToken()
		if err != nil {
			return nil, err
		}

		baseDenomList = others
	} else if strings.HasPrefix(baseDenom, constant.IBCTokenPrefix) {
		var err error
		baseDenomList, err = svc.getBaseOfIBCToken(baseDenom)
		if err != nil {
			return nil, err
		}

		baseDenomList = append(baseDenomList, baseDenom)
	} else {
		baseDenomList = []string{baseDenom}
	}

	return baseDenomList, nil
}

func (svc *TokenService) getBaseOfIBCToken(denom string) ([]string, error) {
	denomList, err := denomRepo.FindByDenom(denom)
	if err != nil {
		return nil, err
	}

	set := utils.NewStringSet()
	for _, v := range denomList {
		set.Add(v.BaseDenom)
	}

	return set.ToSlice(), nil
}

func (svc *TokenService) IBCTokenList(req *vo.IBCTokenListReq) (*vo.IBCTokenListResp, errors.Error) {
	list, err := tokenStatisticsRepo.List(req)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.IBCTokenItem, 0, len(list))
	for _, v := range list {
		if v.DenomAmount != constant.ZeroDenomAmount {
			items = append(items, vo.IBCTokenItem{
				Denom:      v.Denom,
				DenomPath:  v.DenomPath,
				Chain:      v.Chain,
				TokenType:  v.Type,
				IBCHops:    v.IBCHops,
				Amount:     v.DenomAmount,
				ReceiveTxs: v.ReceiveTxs,
			})
		}
	}

	return &vo.IBCTokenListResp{
		Items: items,
	}, nil
}

func (svc *TokenService) IBCTokenListCount(req *vo.IBCTokenListReq) (int64, errors.Error) {
	totalItem, err := tokenStatisticsRepo.CountList(req)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return totalItem, nil
}
