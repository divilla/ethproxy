package application

import (
	"github.com/divilla/ethproxy/interfaces"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	controller struct {
		service *service
		logger  interfaces.Logger
	}
)

func Controller(e *echo.Echo, client interfaces.EthereumHttpClient, cache interfaces.BlockCacher) {
	c := &controller{
		service: Service(client, cache, e.Logger),
	}

	e.GET("/cache-free-space", c.cacheFreeSpace)
	e.GET("/latest-block-number", c.latestBlockNumber)
	e.GET("/block/:bnr", c.getBlockByNumber)
	e.GET("/block/:bnr/transaction/:tid", c.getTransactionByBlockNumberAndIndex)
}

func (c *controller) cacheFreeSpace(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	return ctx.String(http.StatusOK, c.service.cacheFreeSpace())
}

func (c *controller) latestBlockNumber(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	return ctx.String(http.StatusOK, c.service.latestBlockNumber())
}

func (c *controller) getBlockByNumber(ctx echo.Context) error {
	json, err := c.service.getBlockByNumber(ctx.Param("bnr"))
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("Content-Type", "application/json")
	_, err = ctx.Response().Write(json)

	return err
}

func (c *controller) getTransactionByBlockNumberAndIndex(ctx echo.Context) error {
	json, err := c.service.getTransactionByBlockNumberAndIndex(ctx.Param("bnr"), ctx.Param("tid"))
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("Content-Type", "application/json")
	_, err = ctx.Response().Write(json)

	return err
}
