// template.go use template method pattern
// 1. find from redis
// 2. if you not find from redis, get from db then set to redis
// 3. if you find from redis, return

package redis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

// StringTemplate for string
func (r *Client) StringTemplate(key string, expiration time.Duration, doFn func() (interface{}, error)) (string, error) {
	value, err := r.Get(key)
	if err != nil {
		res, err := doFn()
		if err != nil {
			return "", err
		}

		value = formatAny(reflect.ValueOf(res))
		// ignore err when redis set error
		_ = r.Set(key, value, expiration)
		return value, nil
	}

	return value, nil
}

// StringTemplateUnmarshal for string, stores the result in the value pointed to by v
func (r *Client) StringTemplateUnmarshal(key string, expiration time.Duration, doFn func() (interface{}, error), v interface{}) error {
	res, err := r.StringTemplate(key, expiration, doFn)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(res), v); err != nil {
		return err
	}
	return nil
}

func formatAny(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v)
	case reflect.String:
		return v.String()
	case reflect.Slice, reflect.Map, reflect.Struct:
		return utils.MustMarshalJsonToStr(v.Interface())
	default:
		panic(fmt.Errorf("stringTemplate not handle this Kind: %s", v.String()))
	}
}
