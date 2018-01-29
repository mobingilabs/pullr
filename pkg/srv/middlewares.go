package srv

import (
	"time"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

func ElapsedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cid := uuid.NewV4().String()
			c.Set("contextid", cid)
			c.Set("starttime", time.Now())

			// Helper func to print the elapsed time since this middleware. Good to call at end of
			// request handlers, right before/after replying to caller.
			c.Set("fnelapsed", func(ctx echo.Context) {
				start := ctx.Get("starttime").(time.Time)
				glog.Infof("<-- %v, delta: %v", ctx.Get("contextid"), time.Now().Sub(start))
			})

			glog.Infof("--> %v", cid)
			return next(c)
		}
	}
}

func ServerHeaderMiddleware(srvName string, version string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "mobingi:pullr:apiserver:"+version)
			return next(c)
		}
	}
}
