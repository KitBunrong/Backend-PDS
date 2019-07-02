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

func GetAllUser(c echo.Context) error {
	var userList []model.User
	m := map[string]interface{}{}

	query := "SELECT id, phone, fullname, meta FROM pdskeyspace.user"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userList = append(userList, model.User{
			ID:       m["id"].(gocql.UUID),
			Fullname: m["fullname"].(string),
			Phone:    m["phone"].(int),
			Meta:     m["meta"].(map[string]string),
		})
		m = map[string]interface{}{}
	}

	if userList == nil {
		return c.String(http.StatusNotFound, "There is empty data")
	}

	return c.JSON(http.StatusOK, userList)
}

func GetOneUser(c echo.Context) error {
	var user model.User

	id := c.Param("id")

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Requestd ID")
	}

	m := map[string]interface{}{}
	query := "SELECT  id,fullname,phone, meta FROM pdskeyspace.user WHERE id=? LIMIT 1"
	iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		user = model.User{
			ID:       m["id"].(gocql.UUID),
			Fullname: m["fullname"].(string),
			Phone:    m["phone"].(int),
			Meta:     m["meta"].(map[string]string),
		}
	}
	return c.JSON(http.StatusOK, user)
}

func AddUser(c echo.Context) (err error) {
	user := new(model.User)

	var gocqlUUID gocql.UUID
	gocqlUUID = gocql.TimeUUID()

	if err = c.Bind(user); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(user)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	fmt.Println("Creating new user")
	if err := cassandra.Session.Query(`
		INSERT INTO pdskeyspace.user(id, fullname, phone, meta) VALUES (?, ?, ?, ?)`,
		gocqlUUID, user.Fullname, user.Phone, user.Meta).Exec(); err != nil {
		log.Fatal(err)
	}

	user.ID = gocqlUUID
	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) (err error) {
	type DeleteUser struct {
		ID gocql.UUID `json:"id" validate:"required"`
	}

	user := new(DeleteUser)

	if err = c.Bind(user); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(user)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "DELETE FROM pdskeyspace.user WHERE id=?;"
	if err := cassandra.Session.Query(query, user.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "User deleted")
}

func ReplaceUser(c echo.Context) (err error) {
	user := new(model.User)

	if err = c.Bind(user); err != nil {
		return
	}

	validate := validator.New()
	errs := validate.Struct(user)

	if errs != nil {
		fmt.Println(errs.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	query := "UPDATE pdskeyspace.user SET fullname = ?, phone = ?, meta = ? WHERE id = ?;"
	if err := cassandra.Session.Query(query, user.Fullname, user.Phone, user.Meta, user.ID).Exec(); err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "User updated")
}
