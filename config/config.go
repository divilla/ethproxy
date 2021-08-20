package config

import "time"

const (
	ServerAddress      = ":8080"
	EthereumJsonRPCUrl = "https://cloudflare-eth.com"
	CacheCapacity      = 5000
	CacheDefaultTTL    = 5 * time.Second
	CacheRemoveExpired = 3 * time.Second
	LatestBlockRefresh = 1 * time.Second
	FetchRetries       = 3
)
