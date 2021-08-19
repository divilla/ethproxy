package main

import (
	"context"
	"github.com/divilla/ethproxy/config"
	"github.com/divilla/ethproxy/internal/application"
	"github.com/divilla/ethproxy/internal/healthcheck"
	"github.com/divilla/ethproxy/internal/test"
	"github.com/divilla/ethproxy/pkg/blockcache"
	"github.com/divilla/ethproxy/pkg/cmiddleware"
	"github.com/divilla/ethproxy/pkg/ethclient"
	"github.com/divilla/ethproxy/pkg/jsonclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// CPUProfile enables cpu profiling. Note: Default is CPU
	//defer profile.Start(profile.MemProfileHeap, profile.ProfilePath("/home/vito/go/projects/ethproxy/cmd/profile/")).Stop()

	e := echo.New()
	e.HTTPErrorHandler = cmiddleware.HTTPErrorHandler
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
	//e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	//	//Skipper:      middleware.DefaultSkipper,
	//	ErrorMessage: "request timeout, please try again",
	//	Timeout:      3*time.Second,
	//}))

	jClient := jsonclient.New(e.Logger)
	err := jClient.Url(config.EthereumJsonRPCUrl)
	if err != nil {
		panic(err)
	}

	client := ethclient.New(jClient, config.LatestBlockRefreshInterval, e.Logger)
	cache := blockcache.New(config.CacheCapacity, config.CacheRemoveExpired, e.Logger)
	defer func() {
		client.Done()
		cache.Done()
	}()

	application.Controller(e, client, cache)
	healthcheck.Controller(e)
	test.Controller(e)

	go func() {
		if err := e.Start(config.ServerAddress); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
