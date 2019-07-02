package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func DeviceDeviceGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllDevices)
	g.GET("/:id", handlers.GetOneDevice)

	g.POST("/", handlers.AddDevice)
	g.DELETE("/", handlers.DeleteDevice)
	g.PUT("/", handlers.ReplaceDevice)
}
