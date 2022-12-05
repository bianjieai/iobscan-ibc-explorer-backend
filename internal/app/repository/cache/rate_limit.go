package cache

import (
	"context"
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

type RateLimitRepo struct {
	dbr repository.OpenApiKeyRepo
}

var rateLimitScript = `
-- 获取调用脚本时传入的第一个key值（用作限流的 key）
local key = KEYS[1]
-- 获取调用脚本时传入的第一个参数值（限流大小）
local limit = tonumber(ARGV[1])
-- 获取调用脚本时传入的第二个参数值（限流周期）
local expire = tonumber(ARGV[2])

-- 获取当前流量大小
local curentLimit = tonumber(redis.call('get', key) or "0")

-- 是否超出限流
if curentLimit + 1 > limit then
    -- 返回(拒绝)
    return 0
else
    -- 没有超出 value + 1
    local incr = tonumber(redis.call("INCRBY", key, 1))
    if incr == 1 then
        -- 设置过期时间
        redis.call("EXPIRE", key, expire)
    end
    -- 返回(放行)
    return 1
end`

const (
	ratePass = 1
)

func (repo *RateLimitRepo) RateLimit(key string, frequency, cycleTime int) (bool, error) {
	keys := []string{fmt.Sprintf(rateLimit, key)}
	val, err := rc.EvalInt(context.Background(), rateLimitScript, keys, frequency, cycleTime)
	if err != nil {
		return false, err
	}

	if val == ratePass {
		return true, nil
	}

	return false, nil
}
