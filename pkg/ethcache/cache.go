package ethcache

import (
	"fmt"
	"github.com/divilla/ethproxy/config"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type (
	EthereumBlockCache struct {
		items map[uint64]*item
		cap   int
		rwm   sync.RWMutex
		done  chan struct{}
	}

	item struct {
		nr      uint64
		json    []byte
		expires int64
	}
)

//New creates new string EthereumBlockCache
func New(capacity int) *EthereumBlockCache {
	c := &EthereumBlockCache{
		items: make(map[uint64]*item),
		cap:   capacity,
		done:  make(chan struct{}),
	}

	//goroutine that deletes expired items from cache
	go func(c *EthereumBlockCache) {
		for {
			select {
			case <-c.done:
				return
			default:
				time.Sleep(config.RemoveExpiredInterval)
				c.clear()
			}
		}
	}(c)

	return c
}

//Get returns ethereum block json
func (c *EthereumBlockCache) Get(nr uint64) ([]byte, error) {
	if val, ok := c.items[nr]; ok && val.expires < time.Now().UnixNano() {
		err := c.Remove(nr)
		return nil, fmt.Errorf("block expired: %w", err)
	}

	c.rwm.RLock()
	defer c.rwm.RUnlock()

	if val, ok := c.items[nr]; ok {
		return val.json, nil
	}

	return nil, errors.New("block not found")
}

//Put caches ethereum block json
func (c *EthereumBlockCache) Put(nr uint64, json []byte, latestBlockNumber uint64) error {
	i := &item{
		nr:      nr,
		json:    json,
		expires: expires(nr, latestBlockNumber).UnixNano(),
	}

	if _, ok := c.items[nr]; ok {
		return errors.New("block already exists in cache")
	}

	c.rwm.Lock()
	defer c.rwm.Unlock()

	c.clearOne()
	c.items[nr] = i

	return nil
}

func (c *EthereumBlockCache) Remove(nr uint64) error {
	c.rwm.Lock()
	defer c.rwm.Unlock()

	if _, ok := c.items[nr]; ok {
		delete(c.items, nr)
		return nil
	}

	return errors.New("block doesn't exist in cache")
}

func (c *EthereumBlockCache) FreeSpace() int {
	return c.cap - len(c.items)
}

//Done disposes object
func (c *EthereumBlockCache) Done() {
	c.done <- struct{}{}
	close(c.done)
}

func (c *EthereumBlockCache) clear() {
	var items expiries
	var expired []uint64

	c.rwm.Lock()
	defer c.rwm.Unlock()

	l := len(c.items)
	if l == 0 {
		return
	}

	now := time.Now().UnixNano()
	for nr, it := range c.items {
		if it.expires < now {
			expired = append(expired, nr)
		} else {
			items = append(items, it)
		}
	}

	if len(expired) == 0 {
		return
	}

	for _, v := range expired {
		delete(c.items, v)
	}
}

func (c *EthereumBlockCache) clearOne() {
	if len(c.items) < c.cap {
		return
	}

	var firstItem *item
	for _, it := range c.items {
		if firstItem == nil || it.expires < firstItem.expires {
			firstItem = it
		}
	}

	if firstItem == nil {
		panic(errors.New("unable to find first item that expires first"))
	}

	delete(c.items, firstItem.nr)
}
