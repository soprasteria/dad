package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Organizations is the controller type
type Organizations struct {
}

// GetAll organizations from database
func (u *Organizations) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	organizations, err := database.Organizations.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all organizations")
	}
	return c.JSON(http.StatusOK, organizations)
}

// Get organization from database
func (u *Organizations) Get(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	organization, err := database.Organizations.FindByID(id)
	if err != nil || organization.ID.Hex() == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("Organization not found %v", id))
	}
	return c.JSON(http.StatusOK, organization)
}

// Delete organization from database
func (u *Organizations) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.Organizations.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while removing organization: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}

// Save creates or update given organization
func (u *Organizations) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get organization from body
	var organization types.Organization
	err := c.Bind(&organization)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted organization is not valid: %v", err))
	}

	if organization.Name == "" {
		err = errors.New("name field cannot be empty")
	}

	exists, err := database.Organizations.Exists(organization.Name)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error checking while checking if the organization already exists: %v", err))
	}

	if exists {
		return c.String(http.StatusConflict, fmt.Sprintf("Received organization already exists"))
	}

	if id != "" {
		// Organization will be updated
		organization.ID = bson.ObjectIdHex(id)
	} else {
		// Organization will be created
		organization.ID = ""
	}

	organizationSaved, err := database.Organizations.Save(organization)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save organization to database: %v", err))
	}

	return c.JSON(http.StatusOK, organizationSaved)
}
