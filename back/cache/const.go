package cache

import "time"

const (
	ChallengesCacheKey  = "chal%d"
	ChallengesCacheTTL  = time.Hour * 18
	LeaderboardCacheKey = "lbd"
	LeaderboardCacheTTL = time.Hour * 6
	GalaUserCacheTTL    = time.Hour * 24
	SolversCacheKey     = "sol%d"
	SolversCacheTTL     = time.Hour * 1
)
