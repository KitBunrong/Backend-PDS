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

func GetAllDevices(c echo.Context) error {
	var deviceList []model.DeviceRegistration
	m := map[string]interface{}{}

	query := "SELECT id, devicename,registerdate,zone FROM pdskeyspace.dedevice"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		deviceList = append(deviceList, model.DeviceRegistration{
			ID:           m["id"].(gocql.UUID),
			DeviceName:   m["devicename"].(string),
			RegisterDate: m["registerdate"].(string),
			Zone:         m["zone"].(string),
		})
		m = map[string]interface{}{}
	}
	if deviceList == nil {
		return c.String(http.StatusNotFound, "There is empty data")
	}
	return c.JSON(http.StatusOK, deviceList)
}

func GetOneDevice(c echo.Context) error {
	var device model.DeviceRegistration

	id := c.Param("id")

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m := map[string]interface{}{}
	query := "SELECT id,devicename,registerdate,zone FROM pdskeyspace.dedevice WHERE id=? LIMIT 1"
	iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		device = model.DeviceRegistration{
			ID:           m["id"].(gocql.UUID),
			DeviceName:   m["devicename"].(string),
			RegisterDate: m["registerdate"].(string),
			Zone:         m["zone"].(string),
		}
	}
	return c.JSON(http.StatusOK, device)
}

func AddDevice(c echo.Context) (err error) {
	device := new(model.DeviceRegistration)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(device); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(device)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new Device")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.dedevice(id, devicename, registerdate, zone) VALUES (?, ?, ?, ?)`,
		gocqlUUID, device.DeviceName, device.RegisterDate, device.Zone).Exec(); err != nil {
		log.Fatal(err)
	}

	device.ID = gocqlUUID
	return c.JSON(http.StatusOK, device)
}

func DeleteDevice(c echo.Context) (err error) {
	type DeleteDeviceRegisteration struct {
		ID gocql.UUID `json:"id" validate:"required"`
	}

	device := new(DeleteDeviceRegisteration)

	if err = c.Bind(device); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(device)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "DELETE FROM pdskeyspace.dedevice WHERE id=?;"
	if err := cassandra.Session.Query(query, device.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Device deleted")
}

func ReplaceDevice(c echo.Context) (err error) {
	device := new(model.DeviceRegistration)

	if err = c.Bind(device); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(device)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "UPDATE pdskeyspace.dedevice SET devicename = ?,registerdate = ?, zone = ? WHERE id = ?;"
	if err := cassandra.Session.Query(query, device.DeviceName, device.RegisterDate, device.Zone, device.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Device updated")
}
