package middleware

import (
	"strconv"
	"testing"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/stretchr/testify/require"
)

func Test_Middleware(t *testing.T) {
	//for {
	apiKey := "bianjie"
	apiSecret := "123456"
	timestamp := time.Now().Unix()
	host := "http://127.0.0.1:8000"
	uri := "/ibc/txs/BDBA021CF4939699208228E227699C98B48E9FF59C3CD517C2D88281131EA3CC?chain=irishubqa"

	signature := calculateSignature(uri, "", apiSecret, timestamp)
	header := map[string]string{
		"X-Api-Key":   apiKey,
		"X-Timestamp": strconv.FormatInt(timestamp, 10),
		"X-Signature": signature,
	}

	httpCode, bz, err := utils.HttpDo("GET", host+uri, nil, header)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bz))
	require.Equal(t, 200, httpCode)
	//}
}
