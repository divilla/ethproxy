package blockcache

import (
	"github.com/divilla/ethproxy/interfaces"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type (
	EthereumBlockCache struct {
		logger        interfaces.Logger
		items         map[uint64]*item
		capacity      int
		removeExpired time.Duration
		rwm           sync.RWMutex
		done          chan struct{}
	}

	item struct {
		nr      uint64
		json    []byte
		expires int64
	}
)

//New creates new string EthereumBlockCache
func New(logger interfaces.Logger, capacity int, removeExpired time.Duration) *EthereumBlockCache {
	c := &EthereumBlockCache{
		items:         make(map[uint64]*item),
		removeExpired: removeExpired,
		capacity:      capacity,
		logger:        logger,
		done:          make(chan struct{}),
	}

	//goroutine that deletes expired items from cache
	go func(c *EthereumBlockCache) {
		for {
			select {
			case <-c.done:
				return
			case <-time.After(c.removeExpired):
				c.clear()
			}
		}
	}(c)

	return c
}

//Get returns ethereum block json
func (c *EthereumBlockCache) Get(nr uint64) ([]byte, error) {
	c.rwm.RLock()
	defer c.rwm.RUnlock()

	if val, ok := c.items[nr]; ok && val.expires < time.Now().UnixNano() {
		return nil, errors.Errorf("block expired: %s", time.Unix(0, val.expires))
	}

	if val, ok := c.items[nr]; ok {
		return val.json, nil
	}

	return nil, errors.New("block not found")
}

//Put caches ethereum block json
func (c *EthereumBlockCache) Put(nr uint64, json []byte, ttl time.Duration) error {
	i := &item{
		nr:      nr,
		json:    json,
		expires: time.Now().Add(ttl).UnixNano(),
	}

	if val, ok := c.items[nr]; ok && val.expires > time.Now().UnixNano() {
		return errors.Errorf("block number '%d' already exists in cache", nr)
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
	return c.capacity - len(c.items)
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
	if len(c.items) < c.capacity {
		return
	}

	var firstItem *item
	for _, it := range c.items {
		if firstItem == nil || it.expires < firstItem.expires {
			firstItem = it
		}
	}

	if firstItem == nil {
		panic(errors.New("unable to find first item that BlockExpires first"))
	}

	delete(c.items, firstItem.nr)
}
