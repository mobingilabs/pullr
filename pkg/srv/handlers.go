package srv

import (
	"net/http"

	"github.com/labstack/echo"
)

func CopyrightHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.String(http.StatusOK, "Copyright (c) Mobingi, 2015-2017. All rights reserved.")
		return nil
	}
}

func VersionHandler(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.String(http.StatusOK, version)
		return nil
	}
}
