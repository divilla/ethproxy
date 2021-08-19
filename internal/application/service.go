package application

import (
	"fmt"
	"github.com/divilla/ethproxy/interfaces"
	"github.com/divilla/ethproxy/pkg/blockcache"
	"github.com/divilla/ethproxy/pkg/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"strconv"
)

type (
	service struct {
		client interfaces.EthereumHttpClient
		cache  interfaces.BlockCacher
		logger interfaces.Logger
	}
)

func Service(client interfaces.EthereumHttpClient, cache interfaces.BlockCacher, logger interfaces.Logger) *service {
	return &service{
		client: client,
		cache:  cache,
		logger: logger,
	}
}

func (s *service) cacheFreeSpace() string {
	json, err := sjson.Set(`{}`, "cacheFreeSpace", s.cache.FreeSpace())
	if err != nil {
		panic(err)
	}

	return json
}

func (s *service) latestBlockNumber() string {
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

		result := gjson.GetBytes(json, "number")
		if result.Exists() && result.String() != "" {
			nri, err := ethclient.HexToUInt(result.String())
			if err != nil {
				return nil, err
			}

			_ = s.cache.Put(nri, json, blockcache.BlockExpires(nri, s.client.LatestBlockNumber()))
		}

		return json, nil
	}

	nri, err := strconv.ParseUint(nrs, 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("block number '%s' is not valid integer", nrs))
	}

	json, err := s.cache.Get(nri)
	if err == nil {
		return json, nil
	}

	json, err = s.client.GetBlockByNumber(nri)
	if err != nil {
		s.logger.Error(err)
	}
	if json == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("block with number '%s' not found", nrs))
	}

	if err = s.cache.Put(nri, json, blockcache.BlockExpires(nri, s.client.LatestBlockNumber())); err != nil {
		s.logger.Error(err)
	}

	return json, err
}

func (s *service) getTransactionByBlockNumberAndIndex(nrs string, trs string) ([]byte, error) {
	json, err := s.getBlockByNumber(nrs)
	if err != nil {
		return nil, err
	}

	tri, err := strconv.ParseUint(trs, 10, 64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("transaction index '%s' is not valid integer", trs))
	}

	var transaction []byte
	trh := ethclient.UIntToHex(tri)
	gjson.GetBytes(json, "transactions").ForEach(func(key, value gjson.Result) bool {
		if value.Get("transactionIndex").String() == trh {
			transaction = []byte(value.Raw)
			return false
		}
		return true
	})

	if transaction == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("transaction with index '%x' not found in block number '%s'", trs, nrs))
	}

	return transaction, nil
}
