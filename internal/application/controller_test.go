package application

import (
	"github.com/divilla/ethproxy/config"
	"github.com/divilla/ethproxy/pkg/blockcache"
	"github.com/divilla/ethproxy/pkg/ethclient"
	"github.com/divilla/ethproxy/pkg/jsonclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

const RandomizeNumbers = 100
const Requests = 1000
const BlockNumber = 12988583

func TestController_GetLatestBlock(t *testing.T) {
	e := echo.New()

	jClient := jsonclient.New(e.Logger)
	err := jClient.Url(config.EthereumJsonRPCUrl)
	if err != nil {
		panic(err)
	}

	client := ethclient.New(jClient, config.LatestBlockRefreshInterval, e.Logger)
	cache := blockcache.New(config.CacheCapacity, config.CacheRemoveExpired, e.Logger)

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
		assert.True(t, len(rec.Body.String()) > 100)
	}
}

func TestController_GetBlockByNumber(t *testing.T) {
	e := echo.New()

	jClient := jsonclient.New(e.Logger)
	err := jClient.Url(config.EthereumJsonRPCUrl)
	if err != nil {
		panic(err)
	}

	client := ethclient.New(jClient, config.LatestBlockRefreshInterval, e.Logger)
	cache := blockcache.New(config.CacheCapacity, config.CacheRemoveExpired, e.Logger)

	req := httptest.NewRequest("", "/block", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("bnr")
	ctx.SetParamValues("12988583")
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	if assert.NoError(t, c.getBlockByNumber(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, rec.Body.Len() > 100)
	}
}

func TestController_GetBlockByNumberHeavyLoad(t *testing.T) {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	jClient := jsonclient.New(e.Logger)
	err := jClient.Url(config.EthereumJsonRPCUrl)
	if err != nil {
		panic(err)
	}

	client := ethclient.New(jClient, config.LatestBlockRefreshInterval, e.Logger)
	cache := blockcache.New(config.CacheCapacity, config.CacheRemoveExpired, e.Logger)

	req := httptest.NewRequest("", "/block", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	blockRequestMap := make(map[int]int)
	for i := 0; i < Requests; i++ {
		start := time.Now()
		blockNumber := BlockNumber - rand.Intn(RandomizeNumbers)
		blockRequestMap[blockNumber]++

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("bnr")
		ctx.SetParamValues(strconv.Itoa(blockNumber))
		c := &controller{
			service: Service(client, cache, e.Logger),
		}

		if assert.NoError(t, c.getBlockByNumber(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.True(t, rec.Body.Len() > 100)

			stop := time.Now()
			l := log.JSON{}
			l["block_number"] = blockNumber
			l["requests_per_nr"] = blockRequestMap[blockNumber]
			l["cache_space"] = cache.FreeSpace()
			l["latency_human"] = stop.Sub(start).String()
			l["content_length"] = rec.Result().ContentLength
			e.Logger.Infoj(l)
		}
	}
}
