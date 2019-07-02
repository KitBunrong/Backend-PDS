package handlers

import (
	"fmt"
	"log"
	"net/http"

	cassandra "github.com/KitBunrong/pdspart/database"
	"github.com/KitBunrong/pdspart/model"
	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func GetAllCommitlog(c echo.Context) error {
	var commitlog []model.Commitlog
	m := map[string]interface{}{}

	query := "SELECT id, date, meta FROM pdskeyspace.commitlog"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		commitlog = append(commitlog, model.Commitlog{
			ID:   m["id"].(gocql.UUID),
			Date: m["date"].(string),
			Meta: m["meta"].(map[string]string),
		})
		m = map[string]interface{}{}
	}

	if commitlog == nil {
		return c.String(http.StatusNotFound, "There is empty data")
	}

	return c.JSON(http.StatusOK, commitlog)
}

func AddCommitlog(c echo.Context) (err error) {
	commitlog := new(model.Commitlog)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(commitlog); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(commitlog)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new commitlog")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.commitlog(id, date, meta) VALUES (?, ?, ?)`,
		gocqlUUID, commitlog.Date, commitlog.Meta).Exec(); err != nil {
		log.Fatal(err)
	}

	commitlog.ID = gocqlUUID
	return c.JSON(http.StatusOK, commitlog)
}
