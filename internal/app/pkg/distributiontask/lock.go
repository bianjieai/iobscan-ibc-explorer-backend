package distributiontask

import (
	"context"
	"fmt"
	v8 "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

// UnlockScript use lua for redis use
const UnlockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`

const ExpireScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		return 0
	end`

type LockHasExpired interface {
	Lock(key string, value interface{}, expiration time.Duration) error
	UnLock(key string, value interface{}) (interface{}, error)
	ScriptExpire(key string, value interface{}, expiration time.Duration) (interface{}, error)
	TTL(key string) (time.Duration, error)
}

type WithoutLock struct {
}

func (d *WithoutLock) Lock(string, interface{}, time.Duration) error {
	return nil
}

func (d *WithoutLock) UnLock(string, interface{}) (interface{}, error) {
	return nil, nil
}

func (d *WithoutLock) ScriptExpire(key string, value interface{}, expiration time.Duration) (interface{}, error) {
	return true, nil
}

func (d *WithoutLock) TTL(string) (time.Duration, error) {
	return 0, nil
}

type RedisLock struct {
	redisClient v8.UniversalClient
}

func (r *RedisLock) UnLock(key string, value interface{}) (interface{}, error) {
	res, err := r.redisClient.Eval(context.Background(), UnlockScript, []string{key}, value).Result()
	if err != nil {
		logrus.Error("redis execute unlock script fail", err.Error())
	}
	return res, err
}

func (r *RedisLock) ScriptExpire(key string, value interface{}, expiration time.Duration) (interface{}, error) {
	res, err := r.redisClient.Eval(context.Background(), ExpireScript, []string{key}, value, int(expiration/time.Millisecond)).Result()
	if err != nil {
		logrus.Error("redis execute expire script fail, ", err.Error())
	}
	return res, err
}

func (r *RedisLock) TTL(key string) (time.Duration, error) {
	keys := r.redisClient.TTL(context.Background(), key)
	return keys.Result()
}

func (r *RedisLock) Lock(key string, value interface{}, expiration time.Duration) error {
	ok, err := r.redisClient.SetNX(context.Background(), key, value, expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("lock failed, key already use")
	}
	return nil
}
