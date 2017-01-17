package controllers

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Users is the controller type
type Users struct {
}

//GetAll users from database
func (u *Users) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	users, err := database.Users.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all users")
	}
	return c.JSON(http.StatusOK, users)
}

//Get user from database
func (u *Users) Get(c echo.Context) error {
	user := c.Get("user").(types.User)
	return c.JSON(http.StatusOK, user)
}

//Delete user from database
func (u *Users) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.Users.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while remove user: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}
