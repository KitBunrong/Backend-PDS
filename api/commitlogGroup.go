package api

import (
	"github.com/KitBunrong/pdspart/api/handlers"
	"github.com/labstack/echo"
)

func CommitlogGroup(g *echo.Group) {
	g.GET("/", handlers.GetAllCommitlog)
	g.POST("/", handlers.AddCommitlog)
}
