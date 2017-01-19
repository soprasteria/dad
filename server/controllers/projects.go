package controllers

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// Projects is the controller type
type Projects struct {
}

// GetAll functionnal services from database
func (u *Projects) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username": authUser.Username,
		"role":     authUser.Role,
	}).Info("User trying to retrieve all projects")

	var projects []types.Project
	var err error

	switch authUser.Role {
	case types.AdminRole:
		projects, err = database.Projects.FindAll()
	case types.RIRole:
		// TODO: RI can see their projects and the projects associated to their entities
	case types.CPRole:
		// TODO: CP can see their projects only
	default:
		err = fmt.Errorf("Invalid role %s for user %s", authUser.Role, authUser.Username)
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retreiving all functionnal services")
	}
	return c.JSON(http.StatusOK, projects)
}

// Get project from database
func (u *Projects) Get(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	id := c.Param("id")

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username":  authUser.Username,
		"role":      authUser.Role,
		"projectID": id,
	}).Info("User trying to retrieve a project")

	project, err := database.Projects.FindByID(id)
	if err != nil || project.ID.Hex() == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("Project not found %v", id))
	}

	switch authUser.Role {
	case types.AdminRole:
		// No verification, Admin can see everything
	case types.RIRole:
		// TODO: RI can see their projects and the projects associated to their entities
	case types.CPRole:
		// TODO: CP can see their projects only
	}

	return c.JSON(http.StatusOK, project)
}

// Delete project from database
func (u *Projects) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username":  authUser.Username,
		"role":      authUser.Role,
		"projectID": id,
	}).Info("User trying to delete a project")

	// TODO: privileges handling

	res, err := database.Projects.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while removing project: %v", err))
	}

	return c.JSON(http.StatusOK, res)
}

// Save creates or update given project
func (u *Projects) Save(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get project from body
	var project types.Project
	var err error
	err = c.Bind(&project)

	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted project is not valid: %v", err))
	}

	if project.Name == "" || project.Domain == "" {
		err = errors.New("name and domain fields cannot be empty")
	}

	log.WithField("project", project).Info("Received project to save")

	if id != "" {
		// Project will be updated
		project.ID = bson.ObjectIdHex(id)
	} else {
		// Project will be created
		project.ID = ""
	}

	// TODO: add check for Entity and ServiceCenter

	projectSaved, err := database.Projects.Save(project)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save project to database: %v", err))
	}

	return c.JSON(http.StatusOK, projectSaved)
}
