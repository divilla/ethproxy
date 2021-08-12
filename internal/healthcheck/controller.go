package healthcheck

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const Version = "1.0.0"

type (
	controller struct {
		logger echo.Logger
	}
)

func Controller(e *echo.Echo) {
	c := &controller{
		logger: e.Logger,
	}

	e.GET("/healthcheck", c.healthcheck)
	e.HEAD("/healthcheck", c.healthcheck)
}

// healthcheck responds to a healthcheck request.
func (c *controller) healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK "+Version)
}
