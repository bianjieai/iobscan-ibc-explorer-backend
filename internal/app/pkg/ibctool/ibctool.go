package ibctool

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func CalculateDenomValue(priceMap map[string]dto.CoinItem, denom, denomChain string, denomAmount decimal.Decimal) decimal.Decimal {
	key := fmt.Sprintf("%s%s", denom, denomChain)
	coin, ok := priceMap[key]
	if !ok {
		return decimal.Zero
	}

	value := denomAmount.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).
		Mul(decimal.NewFromFloat(coin.Price)).Round(4)
	return value
}

// GetRootDenom get root denom by denom path
//   - fullPath full fullPath, eg："transfer/channel-1/uiris", "uatom"
func GetRootDenom(fullPath string) string {
	split := strings.Split(fullPath, "/")
	return split[len(split)-1]
}

func MatchDcInfo(scChain, scPort, scChannel string, allChainMap map[string]*entity.ChainConfig) (dcChain, dcPort, dcChannel string) {
	if allChainMap == nil || allChainMap[scChain] == nil {
		return
	}

	for _, ibcInfo := range allChainMap[scChain].IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.PortId == scPort && path.ChannelId == scChannel {
				dcChain = path.Chain
				dcPort = path.Counterparty.PortId
				dcChannel = path.Counterparty.ChannelId
				return
			}
		}
	}

	return
}

// TraceDenom trace denom path, parse denom info
//   - fullDenomPath denom full path，eg："transfer/channel-1/uiris", "uatom"
func TraceDenom(fullDenomPath, chain string, allChainMap map[string]*entity.ChainConfig) *entity.IBCDenom {
	unix := time.Now().Unix()
	denom := CalculateIBCHash(fullDenomPath)
	rootDenom := GetRootDenom(fullDenomPath)
	if !strings.HasPrefix(denom, constant.IBCTokenPrefix) { // base denom
		return &entity.IBCDenom{
			Chain:          chain,
			Denom:          denom,
			PrevDenom:      "",
			PrevChain:      "",
			BaseDenom:      denom,
			BaseDenomChain: chain,
			DenomPath:      "",
			RootDenom:      rootDenom,
			IBCHops:        0,
			IsBaseDenom:    true,
			CreateAt:       unix,
			UpdateAt:       unix,
		}
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("trace denom: %s, chain: %s, full path: %s, error. %v ", denom, chain, fullDenomPath, err)
		}
	}()

	var currentChain string
	var isBaseDenom bool
	currentChain = chain
	pathSplits := strings.Split(fullDenomPath, "/")
	denomPath := strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	var TraceDenomList []*dto.DenomSimpleDTO
	TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
		Denom: denom,
		Chain: chain,
	})

	for {
		if len(pathSplits) <= 1 {
			break
		}

		currentPort, currentChannel := pathSplits[0], pathSplits[1]
		tempPrevChain, tempPrevPort, tempPrevChannel := MatchDcInfo(currentChain, currentPort, currentChannel, allChainMap)
		if tempPrevChain == "" { // trace to end
			break
		} else {
			TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
				Denom: CalculateIBCHash(strings.Join(pathSplits[2:], "/")),
				Chain: tempPrevChain,
			})
		}

		currentChain, currentPort, currentChannel = tempPrevChain, tempPrevPort, tempPrevChannel
		pathSplits = pathSplits[2:]
	}

	var prevDenom, prevChain, baseDenom, baseDenomChain string
	if len(TraceDenomList) == 1 { // denom is base denom
		isBaseDenom = true
		baseDenom = denom
		baseDenomChain = chain
	} else {
		isBaseDenom = false
		prevDenom = TraceDenomList[1].Denom
		prevChain = TraceDenomList[1].Chain
		baseDenom = TraceDenomList[len(TraceDenomList)-1].Denom
		baseDenomChain = TraceDenomList[len(TraceDenomList)-1].Chain
	}

	return &entity.IBCDenom{
		Chain:          chain,
		Denom:          denom,
		PrevDenom:      prevDenom,
		PrevChain:      prevChain,
		BaseDenom:      baseDenom,
		BaseDenomChain: baseDenomChain,
		DenomPath:      denomPath,
		RootDenom:      rootDenom,
		IBCHops:        IBCHops(denomPath),
		IsBaseDenom:    isBaseDenom,
		CreateAt:       unix,
		UpdateAt:       unix,
	}
}

func CalculateIBCHash(fullPath string) string {
	if len(strings.Split(fullPath, "/")) == 1 {
		return fullPath
	}

	hash := utils.Sha256(fullPath)
	return fmt.Sprintf("%s/%s", constant.IBCTokenPrefix, strings.ToUpper(hash))
}

// CalculateNextDenomPath calculate full denom path of next hop.
// return full denom path and cross back identification
func CalculateNextDenomPath(packet model.Packet) (string, bool) {
	prefixSc := fmt.Sprintf("%s/%s/", packet.SourcePort, packet.SourceChannel)
	prefixDc := fmt.Sprintf("%s/%s/", packet.DestinationPort, packet.DestinationChannel)
	denomPath := packet.Data.Denom
	if strings.HasPrefix(denomPath, prefixSc) { // transfer to prev chain
		denomPath = strings.Replace(denomPath, prefixSc, "", 1)
		return denomPath, true
	} else {
		denomPath = fmt.Sprintf("%s%s", prefixDc, denomPath)
		return denomPath, false
	}
}

func IBCHops(denomPath string) int {
	return strings.Count(denomPath, constant.IBCHopsIndex)
}

// SplitFullPath get denom path and root denom from denom path
//   - fullPath full fullPath, eg："transfer/channel-1/uiris", "uatom"
func SplitFullPath(fullPath string) (denomPath, rootDenom string) {
	pathSplits := strings.Split(fullPath, "/")
	denomPath = strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	rootDenom = pathSplits[len(pathSplits)-1]
	return
}
