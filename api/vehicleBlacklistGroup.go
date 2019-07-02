package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func VehicleBlacklistGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllBlacklistVehicles)
	g.GET("/:id", handlers.GetOneBlacklistVehicle)

	g.POST("/", handlers.AddBlacklistVehicle)
	g.DELETE("/", handlers.DeleteBlacklistVehicle)
	g.PUT("/", handlers.ReplaceBlacklistVehicle)
}
