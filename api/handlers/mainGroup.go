package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cassandra "github.com/KitBunrong/pdspart/database"
	"github.com/KitBunrong/pdspart/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func Heartbeat(c echo.Context) error {
	return c.String(http.StatusOK, "Status OK")
}

func Signup(c echo.Context) (err error) {
	admin := new(model.Admin)

	// // Bind
	if err = c.Bind(admin); err != nil {
		return
	}

	var gocqlUUID gocql.UUID
	// generate a unique UUID for this user
	gocqlUUID = gocql.TimeUUID()

	// validator.v9
	validate := validator.New()
	errs := validate.Struct(admin)
	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Err")
	}

	// Save user
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.admin (id, fullname, password, email, phone, tokens) VALUES (?, ?, ?, ?, ?, ?)`,
		gocqlUUID, admin.Fullname, admin.Password, admin.Email, admin.Phone, admin.Tokens).Exec(); err != nil {
		log.Fatal(err)
	}

	admin.ID = gocqlUUID
	return c.JSON(http.StatusCreated, admin)
}

func Login(c echo.Context) (err error) {
	admin := new(model.Admin)

	// Bind
	if err = c.Bind(admin); err != nil {
		return
	}

	// validator.v9
	validate := validator.New()
	errs := validate.Struct(admin)
	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Err")
	}

	// Find user
	if err := cassandra.Session.Query(`SELECT id, email FROM pdskeyspace.admin WHERE fullname = ? and password = ? LIMIT 1 ALLOW FILTERING`,
		admin.Fullname, admin.Password).Consistency(gocql.One).Scan(&admin.ID, &admin.Email); err != nil {
		return c.String(http.StatusNotFound, "User is not found")
		// log.Fatal(err)
	}

	//-----
	// JWT
	//-----

	// Create token
	tokens := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := tokens.Claims.(jwt.MapClaims)
	claims["id"] = admin.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response
	admin.Tokens, err = tokens.SignedString([]byte("secretbutsimple"))
	if err != nil {
		return err
	}
	admin.Password = "" // Don't send password
	return c.JSON(http.StatusOK, admin)
}
