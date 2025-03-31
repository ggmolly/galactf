package cache

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

func ReadCached[T any](key string) (*T, error) {
	var result T
	b, err := RedisDb.Get(RedisCtx, key).Bytes()
	if err == redis.Nil {
		return nil, redis.Nil
	} else if err != nil {
		log.Println("[!] failed to read cache:", err)
		return nil, err
	}
	if err := sonic.ConfigFastest.Unmarshal(b, &result); err != nil {
		log.Println("[!] failed to unmarshal cache:", err)
		return nil, err
	}
	return &result, nil
}
