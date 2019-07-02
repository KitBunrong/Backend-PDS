package model

import "github.com/gocql/gocql"

type (
	VehicleManager struct {
		ID           gocql.UUID `json:"id"`
		VehicleName  string     `jon:"vehiclename" validate:"required"`
		PlateNumber  string     `json:"platenumber" validate:"required"`
		RegisterDate string     `json:"registerdate" validate:"required"`
		Location     string     `json:"location" validate:"required"`
		Roll         string     `json:"roll" validate:"required"`
	}
)
