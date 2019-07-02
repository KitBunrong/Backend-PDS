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

func GetAllBlacklistVehicles(c echo.Context) error {
	var vehicleBlacklist []model.VehicleBlacklist
	m := map[string]interface{}{}

	query := "SELECT id, devicename, platenumber, registerdate, location, roll, reason FROM pdskeyspace.veblacklist"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		vehicleBlacklist = append(vehicleBlacklist, model.VehicleBlacklist{
			ID:           m["id"].(gocql.UUID),
			DeviceName:   m["devicename"].(string),
			PlateNumber:  m["platenumber"].(string),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
			Roll:         m["roll"].(string),
			Reason:       m["reason"].(string),
		})
		m = map[string]interface{}{}
	}
	if vehicleBlacklist == nil {
		return c.String(http.StatusInternalServerError, "There is no data")
	}
	return c.JSON(http.StatusOK, vehicleBlacklist)
}

func GetOneBlacklistVehicle(c echo.Context) error {
	var vehicleBlacklist model.VehicleBlacklist

	id := c.Param("id")

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m := map[string]interface{}{}
	query := "SELECT  id,devicename, platenumber, registerdate, location, roll, reason FROM pdskeyspace.veblacklist WHERE id=? LIMIT 1"
	iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		vehicleBlacklist = model.VehicleBlacklist{
			ID:           m["id"].(gocql.UUID),
			DeviceName:   m["devicename"].(string),
			PlateNumber:  m["platenumber"].(string),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
			Roll:         m["roll"].(string),
			Reason:       m["reason"].(string),
		}
	}
	return c.JSON(http.StatusOK, vehicleBlacklist)
}

func AddBlacklistVehicle(c echo.Context) (err error) {
	vehicleBlacklist := new(model.VehicleBlacklist)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(vehicleBlacklist); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleBlacklist)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new user")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.veblacklist(id, devicename, platenumber, registerdate, location, roll, reason) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		gocqlUUID, vehicleBlacklist.DeviceName, vehicleBlacklist.PlateNumber, vehicleBlacklist.RegisterDate, vehicleBlacklist.Location, vehicleBlacklist.Roll, vehicleBlacklist.Reason).Exec(); err != nil {
		log.Fatal(err)
	}

	vehicleBlacklist.ID = gocqlUUID
	return c.JSON(http.StatusOK, vehicleBlacklist)
}

func DeleteBlacklistVehicle(c echo.Context) (err error) {
	type DeleteVehicleBlacklist struct {
		ID gocql.UUID `json:"id" validate:"required"`
	}

	vehicleBlacklist := new(DeleteVehicleBlacklist)

	if err = c.Bind(vehicleBlacklist); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleBlacklist)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "DELETE FROM pdskeyspace.veblacklist WHERE id=?;"
	if err := cassandra.Session.Query(query, vehicleBlacklist.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "User deleted")
}

func ReplaceBlacklistVehicle(c echo.Context) (err error) {
	vehicleBlacklist := new(model.VehicleBlacklist)

	if err = c.Bind(vehicleBlacklist); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleBlacklist)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "UPDATE pdskeyspace.veblacklist SET devicename = ?, platenumber = ?, registerdate = ?, location = ?, roll = ?, reason = ? WHERE id = ?;"
	if err := cassandra.Session.Query(query, vehicleBlacklist.DeviceName, vehicleBlacklist.PlateNumber, vehicleBlacklist.RegisterDate, vehicleBlacklist.Location, vehicleBlacklist.Roll, vehicleBlacklist.Reason, vehicleBlacklist.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Vehicle Blacklist updated")
}
