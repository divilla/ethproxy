package interfaces

import "time"

type BlockCacher interface {
	Get(nr uint64) ([]byte, error)
	Put(nr uint64, json []byte, ttl time.Duration) error
	Remove(nr uint64) error
	FreeSpace() int
}
