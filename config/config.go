package config

import "time"

const (
	ServerAddress              = ":8080"
	EthereumJsonRPCUrl         = "https://cloudflare-eth.com"
	CacheCapacity              = 5000
	CacheShardCount            = 32
	CacheDefaultTTL            = 5 * time.Second
	LatestBlockRefreshInterval = 3 * time.Second
	CacheRemoveExpired         = 3 * time.Second
	FetchRetries               = 3
)
