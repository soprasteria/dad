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

// Entities is the controller type
type Entities struct {
}

// GetAll entities from database
func (u *Entities) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	entities, err := database.Entities.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all entities")
	}
	return c.JSON(http.StatusOK, entities)
}

// Get entity from database
func (u *Entities) Get(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	entity, err := database.Entities.FindByID(id)
	if err != nil || entity.ID.Hex() == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("Entity not found %v", id))
	}
	return c.JSON(http.StatusOK, entity)
}

// Delete entity from database
func (u *Entities) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	res, err := database.Entities.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while removing entity: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}

// Save creates or update given entity
func (u *Entities) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get entity from body
	var entity types.Entity
	err := c.Bind(&entity)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted entity is not valid: %v", err))
	}

	if entity.Name == "" {
		err = errors.New("name field cannot be empty")
	}

	exists, err := database.Entities.Exists(entity.Name)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error checking while checking if the entity already exists: %v", err))
	}

	if exists {
		return c.String(http.StatusConflict, fmt.Sprintf("Received entity already exists"))
	}

	if id != "" {
		// Entity will be updated
		entity.ID = bson.ObjectIdHex(id)
	} else {
		// Entity will be created
		entity.ID = ""
	}

	entitySaved, err := database.Entities.Save(entity)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save entity to database: %v", err))
	}

	return c.JSON(http.StatusOK, entitySaved)
}
