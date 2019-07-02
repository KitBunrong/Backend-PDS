package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/", handlers.Heartbeat)

	e.POST("/signup", handlers.Signup)
	e.POST("/login", handlers.Login)

}
