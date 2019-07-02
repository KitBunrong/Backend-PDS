package router

import (
	"github.com/KitBunrong/pdspart/api"
	"github.com/KitBunrong/pdspart/api/middlewares"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	// create group
	userGroup := e.Group("/users")
	deviceZoneGroup := e.Group("/device/zones")
	deviceDeviceGroup := e.Group("/device/devices")
	vehicleManagerGroup := e.Group("/vehicle/managers")
	vehicleBlacklistGroup := e.Group("/vehicle/blacklists")
	commitlogGroup := e.Group("/log")

	// main middleware
	e.Use(middleware.Logger())
	middlewares.SetMainMiddlewares(e)

	// set group middleware
	// middlewares.SetUserMiddlewares(userGroup)
	// middlewares.SetDeviceMiddlewares(deviceGroup)
	// middlewares.SetVehicleMiddlewares(vehicleGroup)

	//set main route
	api.MainGroup(e)

	// set group route
	api.UserGroup(userGroup)
	api.DeviceZoneGroup(deviceZoneGroup)
	api.DeviceDeviceGroup(deviceDeviceGroup)
	api.VehicleManagerGroup(vehicleManagerGroup)
	api.VehicleBlacklistGroup(vehicleBlacklistGroup)
	api.CommitlogGroup(commitlogGroup)

	return e
}
