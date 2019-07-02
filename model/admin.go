package model

import "github.com/gocql/gocql"

type (
	Admin struct {
		/*
			You have to check condition for if:
				admin give email so phone is optional and if
				admin give phone so email is optional.
		*/
		ID       gocql.UUID `json:"id"`
		Fullname string     `jon:"fullname" validate:"required"`
		Password string     `json:"password" validate:"required"`
		Email    string     `jon:"email" validate:"email"`
		Phone    int        `json:"phone" validate:"required"`
		Tokens   string     `json:"tokens,omitempty"`
	}
)
