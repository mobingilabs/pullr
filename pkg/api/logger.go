package api

import (
	"time"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// LoggerMiddleware logs incoming requests with their response status along with errors
// other than http errors
func LoggerMiddleware(logger domain.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			elapsed := time.Since(start)

			logger.Infof("%s %s - %d Took %.2fms", c.Request().Method, c.Request().URL.Path, c.Response().Status, float64(elapsed)/1000000.0)

			_, isHttpErr := err.(*echo.HTTPError)
			if !isHttpErr && err != nil {
				logger.Errorf("error: %v", err)
			}

			return nil
		}
	}
}
