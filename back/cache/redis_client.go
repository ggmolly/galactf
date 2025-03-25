package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	HashHexLength = 16
)

var (
	RedisCtx = context.Background()
	RedisDb  *redis.Client
)

// https://gist.github.com/ggmolly/e0dc8da8f1d23bd08b3e4c138b8c28b1
// As proven by the above gist, JSON is preferred over Gob for small objects

func InitRedisClient() {
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("error parsing REDIS_DB: %v", err)
	}
	if os.Getenv("REDIS_HOST") == "" {
		log.Fatal("REDIS_HOST is not set")
	}
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	})
	_, err = RedisDb.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}
}
