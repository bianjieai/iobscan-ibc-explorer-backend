package utils

import (
	"encoding/json"
)

func MustMarshalJson(v interface{}) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

func MustUnmarshalJson(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func MustMarshalJsonToStr(v interface{}) string {
	return string(MustMarshalJson(v))
}

func MustUnmarshalJsonStr(data string, v interface{}) {
	MustUnmarshalJson([]byte(data), v)
}

func MarshalJsonIgnoreErr(v interface{}) []byte {
	bz, _ := json.Marshal(v)
	return bz
}

func UnmarshalJsonIgnoreErr(data []byte, v interface{}) {
	_ = json.Unmarshal(data, v)
}
