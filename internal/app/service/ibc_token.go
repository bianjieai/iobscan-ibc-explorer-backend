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
	IBCTokenList(baseDenom string, req *vo.IBCTokenListReq) (*vo.IBCTokenListResp, errors.Error)
	IBCTokenListCount(baseDenom string, req *vo.IBCTokenListReq) (int64, errors.Error)
}

type TokenService struct {
}

var _ ITokenService = new(TokenService)

func (svc *TokenService) List(req *vo.TokenListReq) (*vo.TokenListResp, errors.Error) {
	baseDenomList, err := svc.analyzeBaseDenom(req.BaseDenom)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	list, err := tokenRepo.List(baseDenomList, req.Chain, req.TokenType, skip, limit)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.TokenItem, 0, len(list))
	for _, v := range list {
		items = append(items, vo.TokenItem{
			BaseDenom:         v.BaseDenom,
			ChainId:           v.ChainId,
			TokenType:         v.Type,
			Supply:            v.Supply,
			Currency:          v.Currency,
			Price:             v.Price,
			ChainsInvolved:    v.ChainsInvolved,
			IBCTransferTxs:    v.TransferTxs,
			IBCTransferAmount: v.TransferAmount,
		})
	}

	var totalItem int64
	if req.UseCount {
		totalItem, err = tokenRepo.CountList(baseDenomList, req.Chain, req.TokenType)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.TokenListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}

func (svc *TokenService) ListCount(req *vo.TokenListReq) (int64, errors.Error) {
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
		others, err := getTokenOthers()
		if err != nil {
			return nil, err
		}

		baseDenomList = others
	} else if strings.HasPrefix(baseDenom, constant.IBCTokenPreFix) {
		var err error
		baseDenomList, err = svc.getBaseOfIBCToken(baseDenom)
		if err != nil {
			return nil, err
		}
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

func (svc *TokenService) IBCTokenList(baseDenom string, req *vo.IBCTokenListReq) (*vo.IBCTokenListResp, errors.Error) {
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	list, err := tokenStatisticsRepo.List(baseDenom, req.Chain, req.TokenType, skip, limit)
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
	if req.UseCount {
		totalItem, err = tokenStatisticsRepo.CountList(baseDenom, req.Chain, req.TokenType)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.IBCTokenListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}

func (svc *TokenService) IBCTokenListCount(baseDenom string, req *vo.IBCTokenListReq) (int64, errors.Error) {
	totalItem, err := tokenStatisticsRepo.CountList(baseDenom, req.Chain, req.TokenType)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return totalItem, nil
}

func getTokenOthers() ([]string, error) {
	allBaseDenom, err := baseDenomRepo.FindAll()
	if err != nil {
		return nil, err
	}

	noSymbolDenomList, err := denomRepo.FindNoSymbolDenoms()
	if err != nil {
		return nil, err
	}

	noSymbolDenomSet := utils.NewStringSet()
	for _, v := range noSymbolDenomList {
		noSymbolDenomSet.Add(v.BaseDenom)
	}

	baseDenomMap := allBaseDenom.ConvertToMap()
	for k, _ := range noSymbolDenomSet {
		_, ok := baseDenomMap[k]
		if ok {
			noSymbolDenomSet.Remove(k)
		}
	}

	return noSymbolDenomSet.ToSlice(), nil
}
