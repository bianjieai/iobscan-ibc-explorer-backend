package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
)

type ITxService interface {
	Query(hash string, req vo.TxReq) (*vo.TxResp, errors.Error)
	FailureStatistics(chain string, startTime, endTime int64) (*vo.FailureStatisticsResp, errors.Error)
}

var _ ITxService = new(TxService)

type TxService struct {
}

func (svc *TxService) getChain(chain string) (*entity.ChainConfig, errors.Error) {
	chainConf, err := chainConfigRepo.FindOne(chain)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.WrapBadRequest(fmt.Sprintf("Chain %s not fount", chain))
		}
		return nil, errors.Wrap(err)
	}

	return chainConf, nil
}

func (svc *TxService) Query(hash string, req vo.TxReq) (*vo.TxResp, errors.Error) {
	chainConf, e := svc.getChain(req.Chain)
	if e != nil {
		return nil, e
	}

	txs, err := txRepo.GetIbcTxsByHash(req.Chain, hash)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	locateTx, packetDTO, e := svc.locateTx(txs, req)
	if e != nil {
		return nil, e
	}

	txDetail := vo.BuildTxDetail(req.Chain, locateTx)

	scChain, dcChain := matchChannel(chainConf, packetDTO.Channel, packetDTO.Port, packetDTO.TxType)
	if scChain == "" {
		return &vo.TxResp{
			TxDetail:      &txDetail,
			Ics20Transfer: nil,
		}, nil
	}

	ics20Transfer, e := svc.ics20Transfer(scChain, dcChain, packetDTO.PacketId)
	if e != nil {
		return nil, e
	} else {
		return &vo.TxResp{
			TxDetail:      &txDetail,
			Ics20Transfer: ics20Transfer,
		}, nil
	}
}

func (svc *TxService) locateTx(txs []*entity.Tx, req vo.TxReq) (*entity.Tx, *dto.MatchTxPacketDTO, errors.Error) {
	if len(txs) == 0 {
		return nil, nil, errors.WrapTxNotFound()
	}

	var matchedTxs []*entity.Tx
	var matchPacketInfo []*dto.MatchTxPacketDTO
	for _, tx := range txs {
		for _, msg := range tx.DocTxMsgs {
			txType := entity.TxType(msg.Type)
			if txType != entity.TxTypeTransfer && txType != entity.TxTypeRecvPacket && txType != entity.TxTypeAckPacket && txType != entity.TxTypeTimeoutPacket {
				continue
			}

			if ok, mdto := matchTxPacket(entity.TxType(msg.Type), msg, req.Channel, req.Port, req.PacketSequence); ok {
				matchedTxs = append(matchedTxs, tx)
				matchPacketInfo = append(matchPacketInfo, mdto)
				break
			}
		}
	}

	if len(matchedTxs) == 0 {
		return nil, nil, errors.WrapTxNotFound()
	}

	if len(matchedTxs) > 1 {
		var tips = make([]string, 0, len(matchedTxs))
		for _, v := range matchPacketInfo {
			tips = append(tips, fmt.Sprintf("%s-%s-%s", v.Channel, v.Port, v.Sequence))
		}
		tipMsg := strings.Join(tips, " or ")
		errMsg := fmt.Sprintf("Tx hash isn't unique.Please supply channel-port-packet_sequence parameters: %s", tipMsg)
		return nil, nil, errors.WrapTxNotUnique(errMsg)
	}

	return matchedTxs[0], matchPacketInfo[0], nil
}

func (svc *TxService) ics20Transfer(scChain, dcChain, packetId string) (*vo.Ics20Transfer, errors.Error) {
	if packetId == "" || scChain == "" {
		return nil, nil
	}

	scTxs, err := txRepo.FindByPacketId(scChain, packetId)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	transferTx, ackTxs, timeoutTxs := svc.parseScChainTxs(scChain, packetId, scTxs)
	if transferTx == nil {
		return nil, nil
	}

	recvPacketTxs := make([]vo.SimpleTxExt, 0)
	if dcChain != "" {
		DcTxs, err := txRepo.FindByPacketId(dcChain, packetId)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		recvPacketTxs = svc.parseDcChainTxs(dcChain, packetId, DcTxs)
	}

	transferMsg := transferTx.Msg.TransferMsg()
	var modelPacket model.Packet
	if len(recvPacketTxs) > 0 {
		modelPacket = recvPacketTxs[0].Msg.RecvPacketMsg().Packet
	} else if len(timeoutTxs) > 0 {
		modelPacket = timeoutTxs[0].Msg.TimeoutPacketMsg().Packet
	} else if len(ackTxs) > 0 {
		modelPacket = ackTxs[0].Msg.AckPacketMsg().Packet
	}

	ibcPacket := vo.BuildIBCPacket(scChain, dcChain, modelPacket)
	if ibcPacket.SourceChannel == "" {
		ibcPacket.SourceChannel = transferMsg.SourceChannel
		ibcPacket.SourcePort = transferMsg.SourcePort
	}

	return &vo.Ics20Transfer{
		Sender:           transferMsg.Sender,
		Receiver:         transferMsg.Receiver,
		IBCPacket:        ibcPacket,
		Token:            transferMsg.Token,
		TransferTx:       transferTx,
		RecvPacketTxs:    recvPacketTxs,
		AckPacketTxs:     ackTxs,
		TimeoutPacketTxs: timeoutTxs,
	}, nil
}

func (svc *TxService) parseScChainTxs(chain, packetId string, txs []*entity.Tx) (transferTx *vo.SimpleTx, ackTxs, timeoutTxs []vo.SimpleTxExt) {
	ackTxs = make([]vo.SimpleTxExt, 0)
	timeoutTxs = make([]vo.SimpleTxExt, 0)
	for _, tx := range txs {
		for msgIndex, msg := range tx.DocTxMsgs {
			if msg.CommonMsg().PacketId != packetId {
				continue
			}

			simpleTx := vo.SimpleTx{
				Chain:  chain,
				TxHash: tx.TxHash,
				Height: tx.Height,
				TxTime: tx.Time,
				Status: tx.Status,
				Msg:    msg,
			}

			switch entity.TxType(msg.Type) {
			case entity.TxTypeTransfer:
				transferTx = &simpleTx
			case entity.TxTypeTimeoutPacket:
				timeoutTxs = append(timeoutTxs, vo.SimpleTxExt{
					SimpleTx:    simpleTx,
					IsEffective: isTimeoutPacketEffective(msgIndex, tx),
				})
			case entity.TxTypeAckPacket:
				ackTxs = append(ackTxs, vo.SimpleTxExt{
					SimpleTx:    simpleTx,
					IsEffective: isAckPacketEffective(msgIndex, msg, tx),
				})
			}
		}
	}

	return
}

func (svc *TxService) parseDcChainTxs(chain, packetId string, txs []*entity.Tx) (recvPacketTxs []vo.SimpleTxExt) {
	recvPacketTxs = make([]vo.SimpleTxExt, 0)
	for _, tx := range txs {
		for msgIndex, msg := range tx.DocTxMsgs {
			if msg.CommonMsg().PacketId != packetId {
				continue
			}

			simpleTx := vo.SimpleTx{
				Chain:  chain,
				TxHash: tx.TxHash,
				Height: tx.Height,
				TxTime: tx.Time,
				Status: tx.Status,
				Msg:    msg,
			}

			if entity.TxType(msg.Type) == entity.TxTypeRecvPacket {
				recvPacketTxs = append(recvPacketTxs, vo.SimpleTxExt{
					SimpleTx:    simpleTx,
					IsEffective: isRecvPacketEffective(msgIndex, tx),
				})
				break
			}
		}
	}

	return
}

func matchTxPacket(txType entity.TxType, msg *model.TxMsg, channel, port, sequence string) (bool, *dto.MatchTxPacketDTO) {
	transferFunc := func() (bool, *dto.MatchTxPacketDTO) {
		transferTxMsg := msg.TransferMsg()
		if channel != "" && transferTxMsg.SourceChannel != channel {
			return false, nil
		}

		if port != "" && transferTxMsg.SourcePort != port {
			return false, nil
		}

		if sequence != "" {
			return false, nil
		}

		return true, &dto.MatchTxPacketDTO{
			TxType:   entity.TxTypeTransfer,
			Channel:  transferTxMsg.SourceChannel,
			Port:     transferTxMsg.SourcePort,
			Sequence: "",
			PacketId: transferTxMsg.PacketId,
		}
	}

	recvFunc := func() (bool, *dto.MatchTxPacketDTO) {
		packetMsg := msg.RecvPacketMsg()
		packet := packetMsg.Packet
		if channel != "" && packet.DestinationChannel != channel {
			return false, nil
		}

		if port != "" && packet.DestinationPort != port {
			return false, nil
		}

		if sequence != "" && strconv.FormatInt(packet.Sequence, 10) != sequence {
			return false, nil
		}

		return true, &dto.MatchTxPacketDTO{
			TxType:   entity.TxTypeRecvPacket,
			Channel:  packet.DestinationChannel,
			Port:     packet.DestinationPort,
			Sequence: strconv.FormatInt(packet.Sequence, 10),
			PacketId: packetMsg.PacketId,
		}
	}

	ackFunc := func() (bool, *dto.MatchTxPacketDTO) {
		packetMsg := msg.AckPacketMsg()
		packet := packetMsg.Packet
		if channel != "" && packet.SourceChannel != channel {
			return false, nil
		}

		if port != "" && packet.SourcePort != port {
			return false, nil
		}

		if sequence != "" && strconv.FormatInt(packet.Sequence, 10) != sequence {
			return false, nil
		}

		return true, &dto.MatchTxPacketDTO{
			TxType:   entity.TxTypeAckPacket,
			Channel:  packet.SourceChannel,
			Port:     packet.SourcePort,
			Sequence: strconv.FormatInt(packet.Sequence, 10),
			PacketId: packetMsg.PacketId,
		}
	}

	timeoutFunc := func() (bool, *dto.MatchTxPacketDTO) {
		packetMsg := msg.TimeoutPacketMsg()
		packet := packetMsg.Packet
		if channel != "" && packet.SourceChannel != channel {
			return false, nil
		}

		if port != "" && packet.SourcePort != port {
			return false, nil
		}

		if sequence != "" && strconv.FormatInt(packet.Sequence, 10) != sequence {
			return false, nil
		}

		return true, &dto.MatchTxPacketDTO{
			TxType:   entity.TxTypeTimeoutPacket,
			Channel:  packet.SourceChannel,
			Port:     packet.SourcePort,
			Sequence: strconv.FormatInt(packet.Sequence, 10),
			PacketId: packetMsg.PacketId,
		}
	}

	switch txType {
	case entity.TxTypeTransfer:
		return transferFunc()
	case entity.TxTypeRecvPacket:
		return recvFunc()
	case entity.TxTypeAckPacket:
		return ackFunc()
	case entity.TxTypeTimeoutPacket:
		return timeoutFunc()
	default:
		return false, nil
	}
}

func matchChannel(chainConf *entity.ChainConfig, channel, port string, txType entity.TxType) (scChain, dcChain string) {
	chain := chainConf.ChainName
	var cpChain string
	for _, v := range chainConf.IbcInfo {
		for _, p := range v.Paths {
			if p.ChannelId == channel && p.PortId == port {
				cpChain = p.Chain
				break
			}
		}
	}

	if txType == entity.TxTypeRecvPacket {
		dcChain = chain
		scChain = cpChain
	} else {
		dcChain = cpChain
		scChain = chain
	}

	return
}

func isTimeoutPacketEffective(msgIndex int, tx *entity.Tx) bool {
	if tx.Status == entity.TxStatusFailed {
		return false
	}

	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "transfer" {
				return true
			}
		}
	}
	return false
}

func isAckPacketEffective(msgIndex int, txMsg *model.TxMsg, tx *entity.Tx) bool {
	if tx.Status == entity.TxStatusFailed {
		return false
	}

	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if strings.Contains(txMsg.AckPacketMsg().Acknowledgement, "error") { // ack error
				if evt.Type == "transfer" {
					return true
				}
			} else {
				if evt.Type == "fungible_token_packet" {
					return true
				}
			}
		}
	}

	return false
}

func isRecvPacketEffective(msgIndex int, tx *entity.Tx) bool {
	if tx.Status == entity.TxStatusFailed {
		return false
	}

	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "write_acknowledgement" {
				for _, attr := range evt.Attributes {
					if attr.Key == "packet_ack" {
						return true
					}
				}
			}
		}
	}

	return false
}

func (svc *TxService) FailureStatistics(chain string, startTime, endTime int64) (*vo.FailureStatisticsResp, errors.Error) {
	if _, e := svc.getChain(chain); e != nil {
		return nil, e
	}

	statistics, err := ibcTxFailLogRepo.FailureStatistics(chain, startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.FailureStatisticsItem, 0, len(statistics))
	for _, v := range statistics {
		items = append(items, vo.FailureStatisticsItem{
			FailureReason:         strings.ReplaceAll(v.Code, "_", " "),
			FailureTransferNumber: v.TxsNum,
		})
	}

	return &vo.FailureStatisticsResp{
		Items: items,
		StatisticCaliber: vo.FailureStatisticCaliber{
			TxTimeStart: startTime,
			TxTimeEnd:   endTime,
		},
	}, nil
}
