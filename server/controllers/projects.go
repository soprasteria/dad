package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

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

	projects, err := database.Projects.FindForUser(authUser)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Error while retrieving all functionnal services")
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

	userProjects, err := database.Projects.FindForUser(authUser)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username))
	}

	if !userProjects.ContainsBsonID(project.ID) {
		return c.String(http.StatusForbidden, fmt.Sprintf("User %s cannot see the project %s", authUser.Username, project.ID))
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

	userProjects, err := database.Projects.FindForUser(authUser)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username))
	}

	if !userProjects.ContainsBsonID(bson.ObjectIdHex(id)) {
		return c.String(http.StatusForbidden, fmt.Sprintf("User %s cannot delete the project %s", authUser.Username, id))
	}

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

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username":  authUser.Username,
		"role":      authUser.Role,
		"projectID": id,
	}).Info("User trying to save a project")

	if id != "" {
		userProjects, err := database.Projects.FindForUser(authUser)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username))
		}

		if !userProjects.ContainsBsonID(bson.ObjectIdHex(id)) {
			return c.String(http.StatusForbidden, fmt.Sprintf("User %s cannot update the project %s", authUser.Username, id))
		}
	}

	// Get project from body
	var project types.Project
	var err error

	err = c.Bind(&project)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Posted project is not valid: %v", err))
	}

	log.WithField("project", project).Info("Received project to save")

	if project.Name == "" || project.Domain == "" {
		return c.String(http.StatusBadRequest, "The name and domain fields cannot be empty")
	}

	// If an business unit is provided, check it exists in the entity collection
	if project.BusinessUnit.Hex() != "" {
		entity, err := database.Entities.FindByIDBson(project.BusinessUnit)
		if err == mgo.ErrNotFound {
			return c.String(http.StatusBadRequest, fmt.Sprintf("The business unit %s does not exist", project.BusinessUnit))
		} else if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve business unit %s from database: %v", project.BusinessUnit, err))
		}

		if entity.Type != types.BusinessUnitType {
			return c.String(http.StatusBadRequest, fmt.Sprintf("The entity %s (%s) is not an business unit but a %s", entity.Name, project.BusinessUnit, entity.Type))
		}
	}

	// If a service center is provided, check it exists in the entity collection
	if project.ServiceCenter.Hex() != "" {
		entity, err := database.Entities.FindByIDBson(project.ServiceCenter)
		if err == mgo.ErrNotFound {
			return c.String(http.StatusBadRequest, fmt.Sprintf("The service center %s does not exist", project.ServiceCenter))
		} else if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve service center %s from database: %v", project.ServiceCenter, err))
		}

		if entity.Type != types.ServiceCenterType {
			return c.String(http.StatusBadRequest, fmt.Sprintf("The entity %s (%s) is not an service center but a %s", entity.Name, project.ServiceCenter, entity.Type))
		}
	}

	// Fill ID, Created and Updated fields
	project.Updated = time.Now()
	if id != "" {
		// Project will be updated
		project.ID = bson.ObjectIdHex(id)
	} else {
		// Project will be created
		project.ID = ""
		project.Created = project.Updated
	}

	projectSaved, err := database.Projects.Save(project)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save project to database: %v", err))
	}

	return c.JSON(http.StatusOK, projectSaved)
}
