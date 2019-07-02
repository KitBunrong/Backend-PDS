package model

import "github.com/gocql/gocql"

type (
	DeviceZone struct {
		ID           gocql.UUID `json:"id"`
		RegisterDate string     `json:"registerdate" validate:"required"`
		Location     string     `json:"location" validate:"required"`
	}
)
