package application

import (
	"github.com/divilla/ethproxy/config"
	"github.com/divilla/ethproxy/pkg/blockcache"
	"github.com/divilla/ethproxy/pkg/ethclient"
	"github.com/divilla/ethproxy/pkg/jsonclient"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const RandomizeNumbers = 100
const Requests = 1000
const BlockNumber = 12988583

var (
	e *echo.Echo
	jClient *jsonclient.JsonHttpClient
	client *ethclient.EthereumHttpClient
	cache *blockcache.EthereumBlockCache
)

func init() {
	e = echo.New()
	jClient = jsonclient.New(e.Logger)
	err := jClient.Url(config.EthereumJsonRPCUrl)
	if err != nil {
		panic(err)
	}

	client = ethclient.New(jClient, e.Logger, config.LatestBlockRefresh)
	cache = blockcache.New(e.Logger, config.CacheCapacity, config.CacheRemoveExpired)
}

func TestController_GetLatestBlock(t *testing.T) {
	req := httptest.NewRequest("", "/block", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("bnr")
	ctx.SetParamValues("latest")
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	if assert.NoError(t, c.getBlockByNumber(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		result := gjson.GetBytes(rec.Body.Bytes(), "number")
		assert.True(t, result.Exists())
		nr, err := ethclient.HexToUInt(result.String())
		assert.NoError(t, err)
		assert.True(t, nr > 100000)
	}
}

func TestController_GetBlockByNumber(t *testing.T) {
	req := httptest.NewRequest("", "/block", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("bnr")
	ctx.SetParamValues(strconv.Itoa(BlockNumber))
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	if assert.NoError(t, c.getBlockByNumber(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		result := gjson.GetBytes(rec.Body.Bytes(), "number")
		assert.True(t, result.Exists())
		nr, err := ethclient.HexToUInt(result.String())
		assert.NoError(t, err)
		assert.True(t, nr == uint64(BlockNumber))
	}
}

func TestController_LatestBlockNumber(t *testing.T) {
	req := httptest.NewRequest("", "/latest-block-number", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	if assert.NoError(t, c.latestBlockNumber(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		result := gjson.GetBytes(rec.Body.Bytes(), "latest_block_number")
		assert.True(t, result.Exists())
		assert.True(t, result.Int() > 100000)
	}
}

func TestController_CacheFreeSpace(t *testing.T) {
	req := httptest.NewRequest("", "/cache-free-space", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	if assert.NoError(t, c.cacheFreeSpace(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		result := gjson.GetBytes(rec.Body.Bytes(), "cache_free_space")
		assert.True(t, result.Exists())
		e.Logger.Info(result.Int())
		assert.True(t, result.Int() == config.CacheCapacity - 2)
	}
}

//func TestController_GetBlockByNumberHeavyLoad(t *testing.T) {
//	e := echo.New()
//	e.Logger.SetLevel(log.INFO)
//
//	jClient := jsonclient.New(e.Logger)
//	err := jClient.Url(config.EthereumJsonRPCUrl)
//	if err != nil {
//		panic(err)
//	}
//
//	client := ethclient.New(jClient, config.LatestBlockRefreshInterval, e.Logger)
//	cache := blockcache.New(config.CacheCapacity, config.CacheRemoveExpired, e.Logger)
//
//	req := httptest.NewRequest("", "/block", nil)
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//
//	blockRequestMap := make(map[int]int)
//	for i := 0; i < Requests; i++ {
//		start := time.Now()
//		blockNumber := BlockNumber - rand.Intn(RandomizeNumbers)
//		blockRequestMap[blockNumber]++
//
//		rec := httptest.NewRecorder()
//		ctx := e.NewContext(req, rec)
//		ctx.SetParamNames("bnr")
//		ctx.SetParamValues(strconv.Itoa(blockNumber))
//		c := &controller{
//			service: Service(client, cache, e.Logger),
//		}
//
//		if assert.NoError(t, c.getBlockByNumber(ctx)) {
//			assert.Equal(t, http.StatusOK, rec.Code)
//			assert.True(t, rec.Body.Len() > 100)
//
//			stop := time.Now()
//			l := log.JSON{}
//			l["block_number"] = blockNumber
//			l["requests_per_nr"] = blockRequestMap[blockNumber]
//			l["cache_space"] = cache.FreeSpace()
//			l["latency_human"] = stop.Sub(start).String()
//			l["content_length"] = rec.Result().ContentLength
//			e.Logger.Infoj(l)
//		}
//	}
//}
