package model

import "github.com/gocql/gocql"

type (
	Commitlog struct {
		ID   gocql.UUID        `json:"id"`
		Date string            `jon:"date" validate:"required"`
		Meta map[string]string `json:"meta"`
	}
)
