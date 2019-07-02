package middlewares

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetMainMiddlewares(e *echo.Echo) {
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte("secretbutsimple"),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login request
			if c.Path() == "/" || c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}))
}
