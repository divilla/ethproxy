package application

import (
	"github.com/divilla/ethproxy/pkg/ethcache"
	"github.com/divilla/ethproxy/pkg/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"strconv"
)

type (
	IEthereumCache interface {
		Get(uint64) ([]byte, error)
		Put(uint64, []byte, uint64) error
		Remove(uint64) error
		FreeSpace() int
	}

	IEthereumClient interface {
		LatestBlockNumber() uint64
		GetLatestBlock() ([]byte, error)
		GetBlockByNumber(uint64) ([]byte, error)
	}

	service struct {
		client    IEthereumClient
		cache     IEthereumCache
		logger    echo.Logger
	}
)

func Service(client *ethclient.EthereumHttpClient, cache *ethcache.EthereumBlockCache, logger echo.Logger) *service {
	return &service{
		client:    client,
		cache:     cache,
		logger:    logger,
	}
}

func (s *service) getLatestBlockNumber() string {
	json, err := sjson.Set(`{}`, "latestBlockNumber", s.client.LatestBlockNumber())
	if err != nil {
		panic(err)
	}

	return json
}

func (s *service) getBlockByNumber(nrs string) ([]byte, error) {
	if nrs == "latest" {
		json, err := s.client.GetLatestBlock()
		if err != nil {
			return nil, err
		}

		nri, err := ethclient.HexToUInt(gjson.GetBytes(json, "number").String())
		if err != nil {
			return nil, err
		}

		s.cache.Put(nri, json, s.client.LatestBlockNumber())

		return json, nil
	}

	nri, err := strconv.ParseUint(nrs, 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "block number is not valid integer")
	}

	json, err := s.cache.Get(nri)
	if err != nil {
		return json, nil
	} else {
		s.logger.Error(err)
	}

	json, err = s.client.GetBlockByNumber(nri)
	if err != nil {
		return nil, err
	}
	if json == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "block with that number not found")
	}

	s.cache.Put(nri, json, s.client.LatestBlockNumber())

	return json, err
}

func (s *service) getTransactionByBlockNumberAndIndex(nrs string, trs string) ([]byte, error) {
	json, err := s.getBlockByNumber(nrs)
	if err != nil {
		return nil, err
	}

	tri, err := strconv.ParseUint(trs, 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "transaction index is not valid integer")
	}

	var transaction []byte
	trh := ethclient.UIntToHex(tri)
	gjson.GetBytes(json,"transactions").ForEach(func(key, value gjson.Result) bool {
		if value.Get("transactionIndex").String() == trh {
			transaction = []byte(value.Raw)
			return false
		}
		return true
	})

	if transaction == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "transaction with that index not found in block")
	}

	return transaction, nil
}
