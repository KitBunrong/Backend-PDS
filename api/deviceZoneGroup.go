package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func DeviceZoneGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllZone)
	g.GET("/:id", handlers.GetOneZone)

	g.POST("/", handlers.AddZone)
	g.DELETE("/", handlers.DeleteZone)
	g.PUT("/", handlers.ReplaceZone)
}
