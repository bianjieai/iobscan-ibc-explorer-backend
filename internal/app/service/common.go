package service

import (
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
)

func PubKeyAlgorithm(pubKeyType string) string {
	if strings.Contains(pubKeyType, constant.ETHSECP256K1) {
		return constant.ETHSECP256K1
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
