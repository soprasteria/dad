package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// FunctionalServices is the controller type
type FunctionalServices struct {
}

// GetAll functional services from database
func (u *FunctionalServices) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	functionalServices, err := database.FunctionalServices.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr("Error while retrieving all functional services"))
	}
	return c.JSON(http.StatusOK, functionalServices)
}

// Get functional service from database
func (u *FunctionalServices) Get(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	functionalService, err := database.FunctionalServices.FindByID(id)
	if err != nil || functionalService.ID.Hex() == "" {
		return c.JSON(http.StatusNotFound, types.NewErr(fmt.Sprintf("Functional service not found %v", id)))
	}
	return c.JSON(http.StatusOK, functionalService)
}

// Delete functional service from database
func (u *FunctionalServices) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.FunctionalServices.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while removing functional service: %v", err)))
	}

	return c.JSON(http.StatusOK, res)
}

// Save creates or update given functional service
func (u *FunctionalServices) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get functional service from body
	var functionalService types.FunctionalService

	err := c.Bind(&functionalService)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted functional service is not valid: %v", err)))
	}

	log.WithField("functionalService", functionalService).Info("Received functional service to save")

	if functionalService.Name == "" || functionalService.Package == "" {
		return c.JSON(http.StatusBadRequest, types.NewErr("The name and package fields cannot be empty"))
	}

	exists, err := database.FunctionalServices.Exists(functionalService.Name, functionalService.Package)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Error while checking if the functional service already exists: %v", err)))
	}
	if exists {
		return c.JSON(http.StatusConflict, types.NewErr(fmt.Sprintf("Received functional service already exists")))
	}

	if id != "" {
		// Functional service will be updated
		functionalService.ID = bson.ObjectIdHex(id)
	} else {
		// Functional service will be created
		functionalService.ID = ""
	}

	functionalServiceSaved, err := database.FunctionalServices.Save(functionalService)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to save functional service to database: %v", err)))
	}

	return c.JSON(http.StatusOK, functionalServiceSaved)
}
