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

func GetAllManageVehicles(c echo.Context) error {
	var vehicleManagerList []model.VehicleManager
	m := map[string]interface{}{}

	query := "SELECT id, vehiclename, platenumber, registerdate, location, roll FROM pdskeyspace.vemanager"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		vehicleManagerList = append(vehicleManagerList, model.VehicleManager{
			ID:           m["id"].(gocql.UUID),
			VehicleName:  m["vehiclename"].(string),
			PlateNumber:  m["platenumber"].(string),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
			Roll:         m["roll"].(string),
		})
		m = map[string]interface{}{}
	}
	if vehicleManagerList == nil {
		return c.String(http.StatusInternalServerError, "There is no data")
	}
	return c.JSON(http.StatusOK, vehicleManagerList)
}

func GetOneManageVehicle(c echo.Context) error {
	var vehicleManager model.VehicleManager

	id := c.Param("id")

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m := map[string]interface{}{}
	query := "SELECT  id, vehiclename, platenumber, registerdate, location, roll FROM pdskeyspace.vemanager WHERE id=? LIMIT 1"
	iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		vehicleManager = model.VehicleManager{
			ID:           m["id"].(gocql.UUID),
			VehicleName:  m["vehiclename"].(string),
			PlateNumber:  m["platenumber"].(string),
			RegisterDate: m["registerdate"].(string),
			Location:     m["location"].(string),
			Roll:         m["roll"].(string),
		}
	}
	return c.JSON(http.StatusOK, vehicleManager)
}

func AddManageVehicle(c echo.Context) (err error) {
	vehicleManager := new(model.VehicleManager)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(vehicleManager); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleManager)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new Vehicle Register")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.vemanager(id, vehiclename, platenumber, registerdate, location, roll) VALUES (?, ?, ?, ?, ?, ?)`,
		gocqlUUID, vehicleManager.VehicleName, vehicleManager.PlateNumber, vehicleManager.RegisterDate, vehicleManager.Location, vehicleManager.Roll).Exec(); err != nil {
		log.Fatal(err)
	}

	vehicleManager.ID = gocqlUUID
	return c.JSON(http.StatusOK, vehicleManager)
}

func DeleteManageVehicle(c echo.Context) (err error) {
	type DeleteVehicleManager struct {
		ID gocql.UUID `json:"id" validate:"required"`
	}

	vehicleManager := new(DeleteVehicleManager)

	if err = c.Bind(vehicleManager); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleManager)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "DELETE FROM pdskeyspace.vemanager WHERE id=?;"
	if err := cassandra.Session.Query(query, vehicleManager.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Vehicle deleted")
}

func ReplaceManageVehicle(c echo.Context) (err error) {
	vehicleManager := new(model.VehicleManager)

	if err = c.Bind(vehicleManager); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(vehicleManager)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "UPDATE pdskeyspace.vemanager SET vehiclename = ?, platenumber = ?, registerdate = ?, location = ?, roll =  ? WHERE id = ?;"
	if err := cassandra.Session.Query(query, vehicleManager.VehicleName, vehicleManager.PlateNumber, vehicleManager.RegisterDate, vehicleManager.Location, vehicleManager.Roll, vehicleManager.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Vehicle Manager updated")
}
