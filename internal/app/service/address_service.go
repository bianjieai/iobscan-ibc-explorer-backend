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
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
)

type IAddressService interface {
	BaseInfo(chain, address string) (*vo.BaseInfoResp, errors.Error)
	TxsList(chain, address string, req *vo.AddressTxsListReq) (*vo.AddressTxsListResp, errors.Error)
	TxsCount(chain, address string) (int64, errors.Error)
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
