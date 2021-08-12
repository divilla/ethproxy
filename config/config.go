package config

import "time"

const (
	ServerAddress              = ":8080"
	EthereumJsonRPCUrl         = "https://cloudflare-eth.com"
	LatestBlockRefreshInterval = time.Second
	CacheCapacity              = 60
	CacheDefaultTTL            = 5 * time.Second
	RemoveExpiredInterval      = 3 * time.Second
	FetchRetries               = 3
)
