package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/redis"
)

var rc *redis.Client

func InitRedisClient(c conf.Redis) *redis.Client {
	rc = redis.New(c.Addrs, c.User, c.Password, string(c.Mode), c.Db)
	return rc
}

func GetRedisClient() *redis.Client {
	return rc
}

func RedisStatus() bool {
	return rc.Ping() == nil
}

// RedisDel Redis `DEL` command
func RedisDel(keys ...string) (int64, error) {
	result, err := rc.Del(keys...)
	return result, err
}
