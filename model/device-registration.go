package model

import "github.com/gocql/gocql"

type (
	DeviceRegistration struct {
		ID           gocql.UUID `json:"id"`
		DeviceName   string     `json:"devicename" validate:"required"`
		RegisterDate string     `json:"registerdate" validate:"required"`
		Zone         string     `json:"zone" validate:"required"`
	}
)
