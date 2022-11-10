package service

import (
	"fmt"
	"math"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/shopspring/decimal"
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
