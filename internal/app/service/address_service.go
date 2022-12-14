package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/lcd"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
)

type IAddressService interface {
	BaseInfo(chain, address string) (*vo.BaseInfoResp, errors.Error)
	TxsList(chain, address string, req *vo.AddressTxsListReq) (*vo.AddressTxsListResp, errors.Error)
	TxsCount(chain, address string) (int64, errors.Error)
	TokenList(chain, address string) (*vo.AddrTokenListResp, errors.Error)
	TxsExport(chain, address string) (string, []byte, errors.Error)
}

type AddressService struct {
}

var _ IAddressService = new(AddressService)

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

func (svc *AddressService) TxsList(chain, address string, req *vo.AddressTxsListReq) (*vo.AddressTxsListResp, errors.Error) {
	chainConfigs, err := chainCfgRepo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	allChainMap := make(map[string]*entity.ChainConfig)
	for _, v := range chainConfigs {
		allChainMap[v.ChainName] = v
	}

	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	txs, err := txRepo.GetAddressTxs(chain, address, skip, limit)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	txsItems := make([]vo.AddressTxItem, 0, len(txs))
	for _, tx := range txs {
		txsItems = append(txsItems, svc.loadAddressTxItem(chain, tx, allChainMap))
	}
	page := vo.BuildPageInfo(int64(len(txsItems)), req.PageNum, req.PageSize)
	return &vo.AddressTxsListResp{
		Txs:      txsItems,
		PageInfo: page,
	}, nil
}

func (svc *AddressService) loadAddressTxItem(chain string, tx *entity.Tx, allChainMap map[string]*entity.ChainConfig) vo.AddressTxItem {
	supportTypes := []string{constant.MsgTypeRecvPacket, constant.MsgTypeAcknowledgement, constant.MsgTypeTimeoutPacket, constant.MsgTypeTransfer}
	var denom, amount, port, sender, receiver, scChain, dcChain, baseDenom, baseDenomChain string
	var feeInfo vo.CommonInfo
	var msgType entity.TxType

	getTxType := func() {
		for _, msg := range tx.DocTxMsgs {
			msgType = entity.TxType(msg.Type)
			switch msgType {
			case entity.TxTypeTransfer:
				tm := msg.TransferMsg()
				denom = tm.Token.Denom
				amount = tm.Token.Amount
				port = tm.SourcePort
				sender = tm.Sender
				receiver = tm.Receiver
				scChain = chain
				dcChain, _, _ = matchDcInfo(chain, tm.SourcePort, tm.SourceChannel, allChainMap)
				ibcDenom, err := denomRepo.FindByDenomChain(denom, scChain)
				if err == nil {
					baseDenom = ibcDenom.BaseDenom
					baseDenomChain = ibcDenom.BaseDenomChain
				}
			case entity.TxTypeRecvPacket:
				tm := msg.PacketDataMsg()
				//denom = tm.Packet.Data.Denom
				amount = fmt.Sprint(tm.Packet.Data.Amount)
				port = tm.Packet.DestinationPort
				sender = tm.Packet.Data.Sender
				receiver = tm.Packet.Data.Receiver
				dcChain = chain
				scChain, _, _ = matchDcInfo(chain, tm.Packet.DestinationPort, tm.Packet.DestinationChannel, allChainMap)
				dcDenomFullPath, _ := calculateNextDenomPath(tm.Packet)
				ibcDenom := traceDenom(dcDenomFullPath, chain, allChainMap)
				baseDenom = ibcDenom.BaseDenom
				baseDenomChain = ibcDenom.BaseDenomChain
				denom = ibcDenom.Denom
			case entity.TxTypeAckPacket, entity.TxTypeTimeoutPacket:
				tm := msg.PacketDataMsg()
				//denom = tm.Packet.Data.Denom
				amount = fmt.Sprint(tm.Packet.Data.Amount)
				port = tm.Packet.SourcePort
				sender = tm.Packet.Data.Sender
				receiver = tm.Packet.Data.Receiver
				scChain = chain
				dcChain, _, _ = matchDcInfo(chain, tm.Packet.SourcePort, tm.Packet.SourceChannel, allChainMap)
				ibcDenom := traceDenom(tm.Packet.Data.Denom, chain, allChainMap)
				baseDenom = ibcDenom.BaseDenom
				baseDenomChain = ibcDenom.BaseDenomChain
				denom = ibcDenom.Denom
			}

			if utils.InArray(supportTypes, string(msgType)) {
				break
			}
		}
	}

	getFee := func() {
		for _, v := range tx.Fee.Amount {
			feeInfo.Amount = v.Amount
			feeInfo.Denom = v.Denom
			feeInfo.DenomChain = chain
		}
	}

	getTxType()
	getFee()

	return vo.AddressTxItem{
		TxHash:   tx.TxHash,
		TxStatus: tx.Status,
		TxType:   msgType,
		Port:     port,
		Sender:   sender,
		Receiver: receiver,
		ScChain:  scChain,
		DcChain:  dcChain,
		DenomInfo: vo.DenomInfo{
			CommonInfo: vo.CommonInfo{
				Denom:      denom,
				Amount:     amount,
				DenomChain: chain,
			},
			BaseDenom:      baseDenom,
			BaseDenomChain: baseDenomChain,
		},
		FeeInfo:    feeInfo,
		TxTime:     tx.Time,
		IbcVersion: constant.ICS20,
	}
}

func (svc *AddressService) TxsCount(chain, address string) (int64, errors.Error) {
	count, err := txRepo.CountAddressTxs(chain, address)
	if err != nil {
		return 0, errors.Wrap(err)
	}
	return count, nil
}

func (svc *AddressService) TxsExport(chain, address string) (string, []byte, errors.Error) {
	if len(address) <= 12 {
		return "", nil, errors.WrapDetail(errors.ErrInvalidParams, "invalid address")
	}

	req := vo.AddressTxsListReq{
		Page: vo.Page{
			PageNum:  1,
			PageSize: constant.ExportTxsNum,
		},
		UseCount: false,
	}
	txs, e := svc.TxsList(chain, address, &req)
	if e != nil {
		return "", nil, e
	}

	denomList, err := authDenomRepo.FindAll()
	if err != nil {
		return "", nil, errors.Wrap(err)
	}
	denomMap := denomList.ConvertToMap()

	addressPrefix := address[:6]
	addressSuffix := address[len(address)-6:]
	filename := fmt.Sprintf("%s...%s-%s", addressPrefix, addressSuffix, time.Now().Format(constant.DateFormat))

	var contentArr []string
	header := []string{"Tx Result", "TxHash", "Type", "Port", "From", "To", "Transfer Symbol", "Transfer Amount", "Fee Symbol", "Fee Amount", "Time"}
	contentArr = append(contentArr, strings.Join(header, ","))
	for _, v := range txs.Txs {
		txRes := "Failed"
		if v.TxStatus == entity.TxStatusSuccess {
			txRes = "Success"
		}
		timeStr := strconv.FormatInt(v.TxTime, 10)
		symbol := v.DenomInfo.BaseDenom
		feeSymbol := v.FeeInfo.Denom
		if denom, ok := denomMap[fmt.Sprintf("%s%s", v.DenomInfo.BaseDenomChain, v.DenomInfo.BaseDenom)]; ok {
			symbol = denom.Symbol
		}
		if denom, ok := denomMap[fmt.Sprintf("%s%s", v.FeeInfo.DenomChain, v.FeeInfo.Denom)]; ok {
			feeSymbol = denom.Symbol
		}

		item := []string{txRes, v.TxHash, string(v.TxType), v.Port, v.Sender, v.Receiver, symbol, v.DenomInfo.Amount, feeSymbol, v.FeeInfo.Amount, timeStr}
		contentArr = append(contentArr, strings.Join(item, ","))
	}

	content := strings.Join(contentArr, "\n")
	return filename, []byte(content), nil
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
