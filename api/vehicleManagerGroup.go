package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func VehicleManagerGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllManageVehicles)
	g.GET("/:id", handlers.GetOneManageVehicle)

	g.POST("/", handlers.AddManageVehicle)
	g.DELETE("/", handlers.DeleteManageVehicle)
	g.PUT("/", handlers.ReplaceManageVehicle)
}
