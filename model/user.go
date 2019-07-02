package model

import "github.com/gocql/gocql"

type (
	User struct {
		ID       gocql.UUID        `json:"id"`
		Fullname string            `jon:"fullname" validate:"required"`
		Phone    int               `json:"phone" validate:"required"`
		Meta     map[string]string `json:"meta"`
	}
)
