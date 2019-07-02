package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func UserGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllUser)
	g.GET("/:id", handlers.GetOneUser)

	g.POST("/", handlers.AddUser)
	g.DELETE("/", handlers.DeleteUser)
	g.PUT("/", handlers.ReplaceUser)
}
