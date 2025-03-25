package cache

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/bytedance/sonic"
)

var (
	ErrValueNil = errors.New("value is nil")
)

func WriteInterface(key string, value interface{}, ttl time.Duration) error {
	b, err := sonic.ConfigFastest.Marshal(value)

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return ErrValueNil
		}
	}

	if err != nil {
		log.Println("[!] failed to marshal value:", err)
		return err
	}

	if err := RedisDb.SetEx(RedisCtx, key, b, ttl).Err(); err != nil {
		log.Println("[!] failed to cache value:", err)
		return err
	}

	return nil
}
