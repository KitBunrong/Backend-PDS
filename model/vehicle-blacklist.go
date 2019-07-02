package model

import "github.com/gocql/gocql"

type (
	VehicleBlacklist struct {
		ID           gocql.UUID `json:"id"`
		DeviceName   string     `jon:"devicename" validate:"required"`
		PlateNumber  string        `json:"platenumber" validate:"required"`
		RegisterDate string     `json:"registerdate" validate:"required"`
		Location     string     `json:"location" validate:"required"`
		Roll         string     `json:"roll" validate:"required"`
		Reason       string     `json:"reason" validate:"required"`
	}
)
