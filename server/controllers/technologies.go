package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Technologies is the controller type
type Technologies struct {
}

// GetAll technologies from database
func (u *Technologies) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	technologies, err := database.Technologies.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr("Error while retrieving all technologies"))
	}
	return c.JSON(http.StatusOK, technologies)
}

// Save creates a technology
func (u *Technologies) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	// Get technology from body
	var technology types.Technology

	err := c.Bind(&technology)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted technology is not valid: %v", err)))
	}

	log.WithField("technology", technology).Info("Received technology to save")

	if technology.Name == "" {
		return c.JSON(http.StatusBadRequest, types.NewErr("The name field cannot be empty"))
	}

	exists, err := database.Technologies.Exists(technology.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Error while checking if the technology already exists: %v", err)))
	}
	if exists {
		return c.JSON(http.StatusConflict, types.NewErr(fmt.Sprintf("Received technology already exists")))
	}

	// Technology will be created
	technology.ID = ""

	technologySaved, err := database.Technologies.Save(technology)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to save technology to database: %v", err)))
	}

	return c.JSON(http.StatusOK, technologySaved)
}
