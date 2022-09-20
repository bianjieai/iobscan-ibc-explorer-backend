package enum

type (
	RedisMode string
)

const (
	RedisSingle  RedisMode = "single"  // redis single server
	RedisCluster RedisMode = "cluster" // redis cluster
)
