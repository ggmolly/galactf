package cache

func InvalidateKey(key string) {
	RedisDb.Del(RedisCtx, key)
}
