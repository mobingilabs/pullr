package srv

import (
	"net/http"

	"github.com/labstack/echo"
)

// CopyrightHandler is an echo request handler to respond with copyright text
func CopyrightHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Copyright (c) Mobingi, 2015-2017. All rights reserved.")
	}
}

// VersionHandler is an echo request handle to respond with application version
func VersionHandler(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, version)
	}
}
