package ethclient

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"sync"
	"time"
)

type (
	IHttpClient interface {
		Post(string) ([]byte, error)
	}

	Logger = echo.Logger

	EthereumHttpClient struct {
		client            IHttpClient
		logger            echo.Logger
		refreshInterval   time.Duration
		baseRequest       string
		latestBlockNumber uint64
		done              chan struct{}
		execMap           map[string]exec
	}

	exec struct {
		ch   chan response
		wg   *sync.WaitGroup
	}

	Callback func([]byte, error)
)

func New(client IHttpClient, logger Logger, refreshInterval time.Duration) *EthereumHttpClient {
	c := &EthereumHttpClient{
		client:          client,
		logger:          logger,
		refreshInterval: refreshInterval,
		baseRequest:     `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
		done:            make(chan struct{}),
		execMap:         make(map[string]exec),
	}

	go func(c *EthereumHttpClient) {
		for {
			select {
			case <-c.done:
				return
			default:
				c.setBlockNumber()
				time.Sleep(c.refreshInterval)
			}
		}
	}(c)

	return c
}

func (c *EthereumHttpClient) LatestBlockNumber() uint64 {
	return c.latestBlockNumber
}

func (c *EthereumHttpClient) GetLatestBlock() ([]byte, error) {
	return c.getBlockByNumber("latest")
}

func (c *EthereumHttpClient) GetBlockByNumber(nr uint64) ([]byte, error) {
	return c.getBlockByNumber(UIntToHex(nr))
}

func (c *EthereumHttpClient) Done() {
	//for k := range c.wgMap {
	//	close(c.wgMap[k])
	//}
	c.done <- struct{}{}
	close(c.done)
}

func (c *EthereumHttpClient) setBlockNumber() {
	req := request("blockNumber")
	json, err := c.client.Post(req.String())
	if err != nil {
		c.logger.Error(err)
	}

	resHex, err := parseResponse(json, req)
	if err != nil {
		c.logger.Error(err)
	}

	resInt, err := HexToUInt(string(resHex))
	if err != nil {
		c.logger.Error(fmt.Errorf("unable to parse 'result' from: '%s' to int: %w", json, err))
	}

	c.latestBlockNumber = resInt
}

func (c *EthereumHttpClient) getBlockByNumber(nr string) ([]byte, error) {
	if _, ok := c.execMap[nr]; !ok {
		c.execMap[nr] = exec{
			ch:   make(chan response, 1000),
			wg:   &sync.WaitGroup{},
		}
	}

	// WaitGroup has sole purpose to enable channel drain, I didn't find any other way to detect when last request was issued
	if len(c.execMap[nr].ch) == 0 {
		c.execMap[nr].wg.Add(1)
		defer func() {
			c.execMap[nr].wg.Done()
		}()

		req := request("getBlockByNumber").
			param(nr).
			param(true)

		json, err := c.client.Post(req.String())
		json, err = parseResponse(json, req)

		res := response{
			json: json,
			err:  err,
		}

		//Goroutine is used to drain channel
		go func(c *EthereumHttpClient) {
			time.Sleep(1 * time.Second)
			c.execMap[nr].wg.Wait()
			for len(c.execMap[nr].ch) > 0 {
				<-c.execMap[nr].ch
			}
			return
		}(c)

		//I know it looks ugly, but didn't find any nicer way to push fetch result to unknown number of goroutines
		//i := 0
		for {
			select {
			case c.execMap[nr].ch <- res:
				//i++
				//c.logger.Infof("'%v' request in the same channel", i)
			default:
				return res.json, res.err
			}
		}
	} else {
		c.execMap[nr].wg.Add(1)
		defer func() {
			c.execMap[nr].wg.Done()
		}()

		res := <-c.execMap[nr].ch
		return res.json, res.err
	}
}
