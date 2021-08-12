package ethcache

import (
	"github.com/divilla/ethproxy/config"
	"time"
)

func expires(blockNumber, latestBlockNumber uint64) time.Time {

	// Last 20 blocks have small TTL due to possible reorg
	if blockNumber >= latestBlockNumber - 20 {
		// the further the block from the last one, the more TTL it has.
		return time.Now().Add(config.CacheDefaultTTL)
	}

	// between 20 and 1000 blocks we set TTL depending on distance. The further the block the longer its TTL
	if blockNumber >= latestBlockNumber - 1000 {
		return time.Now().Add(config.CacheDefaultTTL * time.Duration(latestBlockNumber - blockNumber))
	}

	// blocks that are safe to cache get 10 years of TTL
	return time.Now().Add(time.Hour * 24 * 365 * 10)
}
