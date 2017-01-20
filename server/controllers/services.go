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

// FunctionnalServices is the controller type
type FunctionnalServices struct {
}

// GetAll functionnal services from database
func (u *FunctionnalServices) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	functionnalServices, err := database.FunctionnalServices.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all functionnal services")
	}
	return c.JSON(http.StatusOK, functionnalServices)
}

// Get functionnal service from database
func (u *FunctionnalServices) Get(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	functionnalService, err := database.FunctionnalServices.FindByID(id)
	if err != nil || functionnalService.ID.Hex() == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("Functionnal service not found %v", id))
	}
	return c.JSON(http.StatusOK, functionnalService)
}

// Delete functionnal service from database
func (u *FunctionnalServices) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.FunctionnalServices.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while removing functionnal service: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}

// Save creates or update given functionnal service
func (u *FunctionnalServices) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get functionnal service from body
	var functionnalService types.FunctionnalService
	var err error

	err = c.Bind(&functionnalService)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted functionnal service is not valid: %v", err))
	}

	log.WithField("functionnalService", functionnalService).Info("Received functionnal service to save")

	if functionnalService.Name == "" || functionnalService.Package == "" {
		return c.String(http.StatusBadRequest, "The name and package fields cannot be empty")
	}

	exists, err := database.FunctionnalServices.Exists(functionnalService.Name, functionnalService.Package)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error checking while checking if the functionnal service already exists: %v", err))
	}
	if exists {
		return c.String(http.StatusConflict, fmt.Sprintf("Received functionnal service already exists"))
	}

	if id != "" {
		// Functionnal service will be updated
		functionnalService.ID = bson.ObjectIdHex(id)
	} else {
		// Functionnal service will be created
		functionnalService.ID = ""
	}

	functionnalServiceSaved, err := database.FunctionnalServices.Save(functionnalService)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save functionnal service to database: %v", err))
	}

	return c.JSON(http.StatusOK, functionnalServiceSaved)
}
