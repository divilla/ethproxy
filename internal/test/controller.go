package test

import (
	"github.com/divilla/ethproxy/interfaces"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type (
	controller struct {
		logger interfaces.Logger
	}
)

func Controller(e *echo.Echo) {
	c := &controller{
		logger: e.Logger,
	}

	g := e.Group("test")
	g.GET("/panic-recover", c.panicRecover)
	g.GET("/error-recover", c.errorRecover)
	g.GET("/timeout", c.timeout)
	g.GET("/http-error", c.httpError)
}

func (c *controller) panicRecover(ctx echo.Context) error {
	panic(errors.New("panic recover"))
}

func (c *controller) errorRecover(ctx echo.Context) error {
	c.logger.Error(errors.New("some new error"))
	return errors.New("invalid configuration")
}

func (c *controller) timeout(ctx echo.Context) error {
	time.Sleep(6 * time.Second)
	return errors.New("timeout error")
}

func (c *controller) httpError(ctx echo.Context) error {
	return echo.NewHTTPError(http.StatusForbidden, errors.New("timeout error"))
}
