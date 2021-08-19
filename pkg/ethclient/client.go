package ethclient

import (
	"github.com/divilla/ethproxy/interfaces"
	"github.com/labstack/echo/v4"
	"time"
)

type (
	Logger = echo.Logger

	EthereumHttpClient struct {
		client            interfaces.HttpClient
		logger            interfaces.Logger
		refreshInterval   time.Duration
		baseRequest       string
		latestBlockNumber uint64
		done              chan struct{}
		fetchMap          map[string]fetch
	}

	fetch struct {
		add      chan struct{}
		response chan response
	}

	response struct {
		json []byte
		err  error
	}
)

func New(client interfaces.HttpClient, refreshInterval time.Duration, logger interfaces.Logger) *EthereumHttpClient {
	c := &EthereumHttpClient{
		client:          client,
		logger:          logger,
		refreshInterval: refreshInterval,
		baseRequest:     `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
		done:            make(chan struct{}),
		fetchMap:        make(map[string]fetch),
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
	c.done <- struct{}{}
	close(c.done)
}

func (c *EthereumHttpClient) setBlockNumber() {
	req := request("blockNumber")
	json, err := c.client.Post(req.String())
	if err != nil {
		c.logger.Errorf("EthereumHttpClient failed to execute request '%s', with error: %w", req.String(), err)
	}

	resHex, err := parseResponse(json, req)
	if err != nil {
		c.logger.Errorf("EthereumHttpClient failed to parse response '%s' from request '%s' with error: %w", json, req, err)
	}

	resInt, err := HexToUInt(string(resHex))
	if err != nil {
		c.logger.Errorf("EthereumHttpClient failed to parse hex '%s' to int, with error: %w", resHex, err)
	}

	c.latestBlockNumber = resInt
}

func (c *EthereumHttpClient) getBlockByNumber(nr string) ([]byte, error) {
	if _, ok := c.fetchMap[nr]; !ok {
		c.fetchMap[nr] = fetch{
			add:      make(chan struct{}, 1000),
			response: make(chan response),
		}

		go func(c *EthereumHttpClient) {
			req := request("getBlockByNumber").
				param(nr).
				param(true)

			json, err := c.client.Post(req.String())
			json, err = parseResponse(json, req)

			res := response{
				json: json,
				err:  err,
			}

			for {
				select {
				case <- c.fetchMap[nr].add:
					c.fetchMap[nr].response <-res
				case <-time.After(time.Second):
					if len(c.fetchMap[nr].add) == 0 {
						close(c.fetchMap[nr].response)
						delete(c.fetchMap, nr)
						return
					}
				}
			}
		}(c)

		c.fetchMap[nr].add <- struct{}{}

		res := <- c.fetchMap[nr].response
		return res.json, res.err
	} else {
		c.fetchMap[nr].add <- struct{}{}

		res := <- c.fetchMap[nr].response
		return res.json, res.err
	}
}
