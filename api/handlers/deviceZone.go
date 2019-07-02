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

func GetAllZone(c echo.Context) error {
	var zoneList []model.DeviceZone
	m := map[string]interface{}{}

	query := "SELECT id, registerdate, location FROM pdskeyspace.dezone"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		zoneList = append(zoneList, model.DeviceZone{
			ID:           m["id"].(gocql.UUID),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
		})
		m = map[string]interface{}{}
	}
	if zoneList == nil {
		return c.String(http.StatusInternalServerError, "There is no data")
	}
	return c.JSON(http.StatusOK, zoneList)
}

func GetOneZone(c echo.Context) error {
	var zone model.DeviceZone

	id := c.Param("id")

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m := map[string]interface{}{}
	query := "SELECT  id,registerdate,location FROM pdskeyspace.dezone WHERE id=? LIMIT 1"
	iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		zone = model.DeviceZone{
			ID:           m["id"].(gocql.UUID),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
		}
	}
	return c.JSON(http.StatusOK, zone)
}

func AddZone(c echo.Context) (err error) {
	zone := new(model.DeviceZone)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(zone); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(zone)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new user")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.dezone(id, registerdate, location) VALUES (?, ?, ?)`,
		gocqlUUID, zone.RegisterDate, zone.Location).Exec(); err != nil {
		log.Fatal(err)
	}

	zone.ID = gocqlUUID
	return c.JSON(http.StatusOK, zone)
}

func DeleteZone(c echo.Context) (err error) {
	type DeleteDeviceZone struct {
		ID gocql.UUID `json:"id" validate:"required"`
	}

	zone := new(DeleteDeviceZone)

	if err = c.Bind(zone); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(zone)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "DELETE FROM pdskeyspace.dezone WHERE id=?;"
	if err := cassandra.Session.Query(query, zone.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Zone deleted")
}

func ReplaceZone(c echo.Context) (err error) {
	zone := new(model.DeviceZone)

	if err = c.Bind(zone); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(zone)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "UPDATE pdskeyspace.dezone SET registerdate = ?, location = ? WHERE id = ?;"
	if err := cassandra.Session.Query(query, zone.RegisterDate, zone.Location, zone.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Zone updated")
}
