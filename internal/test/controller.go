package test

import (
	"github.com/divilla/ethproxy/interfaces"
	"github.com/divilla/ethproxy/pkg/gerror"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type (
	controller struct {
		logger interfaces.ErrorLogger
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
	panic(gerror.New("panic recover"))
}

func (c *controller) errorRecover(ctx echo.Context) error {
	c.logger.Error(gerror.New("error recover"))
	return gerror.NewCode(gerror.CodeInvalidConfiguration, "invalid configuration")
}

func (c *controller) timeout(ctx echo.Context) error {
	time.Sleep(6 * time.Second)
	return gerror.NewCode(gerror.CodeInvalidConfiguration, "timeout error")
}

func (c *controller) httpError(ctx echo.Context) error {
	return echo.NewHTTPError(http.StatusForbidden, gerror.NewCode(gerror.CodeInvalidConfiguration, "timeout error"))
}
