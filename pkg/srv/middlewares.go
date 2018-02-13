package srv

import (
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// ElapsedMiddleware logs elapsed time while handling the request
func ElapsedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			res := next(c)

			elapsed := time.Since(startTime)
			log.Infof("%s %s %d took %s", c.Request().Method, c.Request().URL.Path, c.Response().Status, elapsed)

			return res
		}
	}
}

// ServerHeaderMiddleware adds server name to response headers
func ServerHeaderMiddleware(srvName string, version string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "mobingi:pullr:apiserver:"+version)
			return next(c)
		}
	}
}
