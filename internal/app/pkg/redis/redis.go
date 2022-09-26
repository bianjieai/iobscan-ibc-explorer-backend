package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	v8 "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const (
	Cluster = "cluster"
	Single  = "single"
)

type Client struct {
	redisClient v8.UniversalClient
}

func New(addrs, user, password, mode string, db int) *Client {
	var redisClient v8.UniversalClient
	if mode == Cluster {
		redisClient = v8.NewClusterClient(&v8.ClusterOptions{
			Addrs:    strings.Split(addrs, ","),
			Username: user,
			Password: password,
		})
	} else if mode == Single {
		redisClient = v8.NewClient(&v8.Options{
			Addr:     addrs,
			Username: user,
			Password: password,
			DB:       db,
		})
	} else {
		logrus.Fatal("unknown redis server mode")
	}
	return &Client{
		redisClient: redisClient,
	}
}

func (r *Client) Close() {
	_ = r.redisClient.Close()
}

// Get RedisClient `GET` command. It returns redis.Nil error when key does not exist
func (r *Client) Get(key string) (string, error) {
	result, err := r.redisClient.Get(context.Background(), key).Result()
	if err != v8.Nil && err != nil {
		logrus.Error("redis get fail, ", err.Error())
	}
	return result, err
}

// Set RedisClient `SET` command.
func (r *Client) Set(key string, value interface{}, expiration time.Duration) error {
	_, err := r.redisClient.Set(context.Background(), key, value, expiration).Result()
	if err != nil {
		logrus.Error("redis Set fail, ", err.Error())
	}
	return err
}

// SMembers RedisClient `SMembers` command. It returns redis.Nil error when key does not exist
func (r *Client) SMembers(key string) ([]string, error) {
	result, err := r.redisClient.SMembers(context.Background(), key).Result()
	if err != v8.Nil && err != nil {
		logrus.Error("redis smembers fail, ", err.Error())
	}
	return result, err
}

// SAdd RedisClient `SAdd` command.
func (r *Client) SAdd(key string, members ...string) (int64, error) {
	n, err := r.redisClient.SAdd(context.Background(), key, members).Result()
	if err != nil {
		logrus.Error("redis hSet fail, ", err.Error())
	}
	return n, err
}

// Del RedisClient `DEL` command
func (r *Client) Del(keys ...string) (int64, error) {
	result, err := r.redisClient.Del(context.Background(), keys...).Result()
	return result, err
}

// HSet accepts values in following formats:
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
//
// Note that it requires RedisClient v4 for multiple field/value pairs support.
func (r *Client) HSet(key string, values ...interface{}) (int64, error) {
	result, err := r.redisClient.HSet(context.Background(), key, values...).Result()
	if err != nil {
		logrus.Error("redis hSet fail, ", err.Error())
	}
	return result, err
}

// HGet RedisClient `HGET` command
func (r *Client) HGet(key, field string) (string, error) {
	result, err := r.redisClient.HGet(context.Background(), key, field).Result()
	if err != nil && err != v8.Nil {
		logrus.Error("redis HGet fail, ", err.Error())
	}
	return result, err
}

// HGetAll RedisClient `HGETALL` command
func (r *Client) HGetAll(key string) (map[string]string, error) {
	result, err := r.redisClient.HGetAll(context.Background(), key).Result()
	if err != nil && err != v8.Nil {
		logrus.Error("redis HGetAll fail, ", err.Error())
	}
	return result, err
}

// UnmarshalGet RedisClient `GET` command with unmarshal. It returns redis.Nil error when key does not exist
func (r *Client) UnmarshalGet(key string, value interface{}) error {
	result, err := r.Get(key)

	if err = json.Unmarshal([]byte(result), value); err != nil {
		return err
	}

	return nil
}

// MarshalSet RedisClient `SET` command with marshal
func (r *Client) MarshalSet(key string, value interface{}, expiration time.Duration) error {
	str, _ := json.Marshal(value)
	return r.Set(key, str, expiration)
}

// UnmarshalSMembers RedisClient `SMembers` command with unmarshal. It returns redis.Nil error when key does not exist
func (r *Client) UnmarshalSMembers(key string, members interface{}) error {
	result, err := r.SMembers(key)
	if err != nil {
		return err
	}

	bz := r.sliceToBytes(result)
	if err = json.Unmarshal(bz, members); err != nil {
		return err
	}
	return nil
}

// MarshalSAdd RedisClient `SAdd` command with marshal
func (r *Client) MarshalSAdd(key string, members ...interface{}) (int64, error) {
	var marshalMembers []interface{}
	for _, m := range members {
		marshal, err := json.Marshal(m)
		if err != nil {
			return 0, err
		}
		marshalMembers = append(marshalMembers, marshal)
	}

	n, err := r.redisClient.SAdd(context.Background(), key, marshalMembers).Result()
	if err != nil {
		logrus.Error("redis hSet fail, ", err.Error())
	}
	return n, err
}

// MarshalHSet accepts values in following formats:
//   - MarshalHSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (r *Client) MarshalHSet(key string, valueMap interface{}) (int64, error) {
	var m map[string]interface{}
	bz, _ := json.Marshal(valueMap)
	if err := json.Unmarshal(bz, &m); err != nil {
		return 0, err
	}

	for k, v := range m {
		marshal, _ := json.Marshal(v)
		m[k] = marshal
	}

	return r.HSet(key, m)
}

// UnmarshalHGet RedisClient `HGET` command with unmarshal
func (r *Client) UnmarshalHGet(key, field string, value interface{}) error {
	result, err := r.HGet(key, field)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(result), value); err != nil {
		return err
	}
	return nil
}

// UnmarshalHGetAll RedisClient `HGETALL` command with unmarshal
// when value is string
func (r *Client) UnmarshalHGetAll(key string, values interface{}) error {
	result, err := r.HGetAll(key)
	if err != nil {
		return err
	}

	bz := r.mapStringToBytes(result)
	if err = json.Unmarshal(bz, values); err != nil {
		return err
	}

	return nil
}

// Expire RedisClient `expire` command
func (r *Client) Expire(key string, expiration time.Duration) bool {
	result, err := r.redisClient.Expire(context.Background(), key, expiration).Result()
	if err != nil && err != v8.Nil {
		logrus.Error("redis Expire fail, ", err.Error())
	}
	return result
}

func (r *Client) sliceToBytes(ss []string) []byte {
	var res bytes.Buffer
	res.WriteString("[")
	for i, v := range ss {
		res.WriteString(v)
		if i != len(ss)-1 {
			res.WriteString(",")
		}
	}
	res.WriteString("]")
	return res.Bytes()
}

func (r *Client) Lock(key string, value interface{}, expiration time.Duration) error {
	ok, err := r.redisClient.SetNX(context.Background(), key, value, expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("lock failed, key already use")
	}
	return nil
}

func (r *Client) Incr(key string) (int64, error) {
	incr, err := r.redisClient.Incr(context.Background(), key).Result()
	if err != nil {
		logrus.Error("redis incr fail, ", err.Error())
	}
	return incr, err
}

func (r *Client) mapStringToBytes(ms map[string]string) []byte {
	var res bytes.Buffer
	res.WriteString("{")
	l := len(ms)
	index := 1
	for k, v := range ms {
		res.WriteString(fmt.Sprintf("\"%s\" : %s", k, v))
		if index != l {
			res.WriteString(",")
		}
		index++
	}
	res.WriteString("}")
	return res.Bytes()
}

func (r *Client) Ping() error {
	err := r.redisClient.Ping(context.Background()).Err()
	if err != nil {
		logrus.Error("redis ping fail, ", err.Error())
	}
	return err
}

// XAdd values accepts values in the following formats:
//   - values = []interface{}{"key1", "value1", "key2", "value2"}
//   - values = []string("key1", "value1", "key2", "value2")
//   - values = map[string]interface{}{"key1": "value1", "key2": "value2"}
func (r *Client) XAdd(stream string, values interface{}) (string, error) {
	messageID, err := r.redisClient.XAdd(context.Background(), &v8.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
	return messageID, err
}

func (r *Client) XDel(stream string, ids ...string) (int64, error) {
	result, err := r.redisClient.XDel(context.Background(), stream, ids...).Result()
	return result, err
}

func (r *Client) XRead(args *v8.XReadArgs) ([]v8.XStream, error) {
	result, err := r.redisClient.XRead(context.Background(), args).Result()
	return result, err
}

func (r *Client) XLen(stream string) (length int64, err error) {
	length, err = r.redisClient.XLen(context.Background(), stream).Result()
	return
}
