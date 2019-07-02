package main

import (
	"fmt"

	cassandra "github.com/KitBunrong/pdspart/database"
	"github.com/KitBunrong/pdspart/router"
)

func main() {
	fmt.Println("Welcome to webserver")

	CassandraSession := cassandra.Session
	defer CassandraSession.Close()

	e := router.New()

	e.Logger.Fatal(e.Start(":8686"))
}
