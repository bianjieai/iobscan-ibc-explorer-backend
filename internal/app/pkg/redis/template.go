// template.go use template method pattern
// 1. find from redis
// 2. if you not find from redis, get from db then set to redis
// 3. if you find from redis, return

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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

// EvalInt execute command `eval`
func (r *Client) EvalInt(ctx context.Context, script string, keys []string, args ...interface{}) (int64, error) {
	val, err := r.redisClient.Eval(ctx, script, keys, args).Result()
	if err != nil {
		logrus.Errorf("redis eval error, %v", err)
		return 0, err
	}

	res, err := formatInt(reflect.ValueOf(val))
	if err != nil {
		logrus.Errorf("redis eval formatInt error, %v", err)
		return 0, err
	}

	return res, nil
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
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Pointer:
		bz, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(bz)
	default:
		panic(fmt.Errorf("stringTemplate not handle this Kind: %s", v.String()))
	}
}

func formatInt(v reflect.Value) (int64, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(v.Float()), nil
	case reflect.String:
		val, err := strconv.Atoi(v.String())
		return int64(val), err
	default:
		return 0, fmt.Errorf("can't convert type %s to int64", v.String())
	}
}
