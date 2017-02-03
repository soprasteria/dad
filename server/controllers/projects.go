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

// GetAll functional services from database
func (u *Projects) GetAll(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)

	authUser := c.Get("authuser").(types.User)
	log.WithFields(log.Fields{
		"username": authUser.Username,
		"role":     authUser.Role,
	}).Info("User trying to retrieve all projects")

	projects, err := database.Projects.FindForUser(authUser)

	if err != nil {
		log.WithError(err).Error("Error while retrieving projects")
		return c.JSON(http.StatusInternalServerError, types.NewErr("Error while retrieving projects"))
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
		return c.JSON(http.StatusNotFound, types.NewErr(fmt.Sprintf("Project not found %v", id)))
	}

	userProjects, err := database.Projects.FindForUser(authUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username)))
	}

	if !userProjects.ContainsBsonID(project.ID) {
		return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf("User %s cannot see the project %s", authUser.Username, project.ID)))
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
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username)))
	}

	if !userProjects.ContainsBsonID(bson.ObjectIdHex(id)) {
		return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf("User %s cannot delete the project %s", authUser.Username, id)))
	}

	res, err := database.Projects.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while removing project: %v", err)))
	}

	return c.JSON(http.StatusOK, res)
}

func canAddEntityToProject(entityToSet, entityFromDB string, authUser types.User) bool {
	if authUser.Role == types.AdminRole {
		return true
	}

	// If the user is a RI, he can only add an entity if:
	// * it's one of his own assigned entities
	// * it's the currently assigned entity of the project
	allowedEntities := make([]bson.ObjectId, len(authUser.Entities))
	copy(allowedEntities, authUser.Entities)
	if bson.IsObjectIdHex(entityFromDB) {
		allowedEntities = append(allowedEntities, bson.ObjectIdHex(entityFromDB))
	}
	for _, allowedEntity := range allowedEntities {
		if bson.ObjectIdHex(entityToSet) == allowedEntity {
			return true
		}
	}
	return false
}

func validateEntity(entityRepo types.EntityRepo, entityToSet, entityFromDB string, entityType types.EntityType, authUser types.User) (int, string) {
	if entityToSet != "" {
		entity, err := entityRepo.FindByID(entityToSet)
		if err == mgo.ErrNotFound {
			return http.StatusBadRequest, fmt.Sprintf("The %s %s does not exist", entityType, entityToSet)
		} else if err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve %s %s from database: %v", entityType, entityToSet, err)
		}

		if !canAddEntityToProject(entityToSet, entityFromDB, authUser) {
			return http.StatusBadRequest, fmt.Sprintf("You can't add the entity %s to a project", entityToSet)
		}

		if entity.Type != entityType {
			return http.StatusBadRequest, fmt.Sprintf("The entity %s (%s) is not of type %s but  %s", entity.Name, entityToSet, entityType, entity.Type)
		}
	}
	return http.StatusOK, ""
}

func validateEntities(entityRepo types.EntityRepo, projectToSave, projectFromDB types.Project, authUser types.User) (int, string) {
	if projectToSave.BusinessUnit == "" && projectToSave.ServiceCenter == "" {
		return http.StatusBadRequest, "At least one of the business unit and service center fields is mandatory"
	}

	// If a business unit is provided, check it exists in the entity collection
	if statusCode, errMessage := validateEntity(entityRepo, projectToSave.BusinessUnit, projectFromDB.BusinessUnit, types.BusinessUnitType, authUser); errMessage != "" {
		return statusCode, errMessage
	}

	// If a service center is provided, check it exists in the entity collection
	if statusCode, errMessage := validateEntity(entityRepo, projectToSave.ServiceCenter, projectFromDB.ServiceCenter, types.ServiceCenterType, authUser); errMessage != "" {
		return statusCode, errMessage
	}

	return http.StatusOK, ""
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

	var projectFromDB types.Project
	if id != "" {
		// Get only projects that the user can modify
		userProjects, err := database.Projects.FindModifiableForUser(authUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while retrieving the projects of the user %s", authUser.Username)))
		}

		if !userProjects.ContainsBsonID(bson.ObjectIdHex(id)) {
			log.WithFields(log.Fields{
				"username":  authUser.Username,
				"role":      authUser.Role,
				"projectID": id,
			}).Info("User isn't allowed to update the project")
			return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf("User %s isn't allowed to update the project", authUser.Username)))
		}

		projectFromDB, err = database.Projects.FindByID(id)
		if err != nil || projectFromDB.ID.Hex() == "" {
			return c.JSON(http.StatusBadRequest, types.NewErr("Trying to modify a non existing project"))
		}

	}

	// Get project from body
	var projectToSave types.Project
	var err error

	err = c.Bind(&projectToSave)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted project is not valid: %v", err)))
	}

	log.WithField("project", projectToSave).Info("Received project to save")

	if projectToSave.Name == "" {
		return c.JSON(http.StatusBadRequest, types.NewErr("The name field cannot be empty"))
	}

	// Not possible to create or update a project with a name already used by another one project
	if existingProject, err := database.Projects.FindByName(projectToSave.Name); err != nil {
		if err != mgo.ErrNotFound {
			return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Can't check whether the project exist in database: %v", err)))
		}
	} else if existingProject.ID != projectToSave.ID {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Another project already exists with the same name %q", existingProject.Name)))
	} else if authUser.Role == types.CPRole {
		// A Project Manager can't update details, if any of the details has changed it's an issue and we shouldn't update the project
		if projectToSave.Domain != existingProject.Domain ||
			projectToSave.ProjectManager != existingProject.ProjectManager ||
			projectToSave.ServiceCenter != existingProject.ServiceCenter ||
			projectToSave.BusinessUnit != existingProject.BusinessUnit {
			log.WithFields(log.Fields{
				"username":                       authUser.Username,
				"role":                           authUser.Role,
				"projectID":                      existingProject.ID,
				"projectToSave.Domain":           projectToSave.Domain,
				"existingProject.Domain":         existingProject.Domain,
				"projectToSave.ProjectManager":   projectToSave.ProjectManager,
				"existingProject.ProjectManager": existingProject.ProjectManager,
				"projectToSave.ServiceCenter":    projectToSave.ServiceCenter,
				"existingProject.ServiceCenter":  existingProject.ServiceCenter,
				"projectToSave.BusinessUnit":     projectToSave.BusinessUnit,
				"existingProject.BusinessUnit":   existingProject.BusinessUnit,
			}).Warn("User isn't allowed to update the project")
			return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("A project manager isn't allowed to update project details")))
		}
	}

	// Check rights to add entities to the project
	httpStatusCode, errorMessage := validateEntities(database.Entities, projectToSave, projectFromDB, authUser)
	if errorMessage != "" {
		return c.JSON(httpStatusCode, types.NewErr(errorMessage))
	}

	// Fill ID, Created and Updated fields
	projectToSave.Updated = time.Now()
	if id != "" {
		// Project will be updated
		projectToSave.ID = bson.ObjectIdHex(id)
	} else {
		// Project will be created
		projectToSave.ID = ""
		projectToSave.Created = projectToSave.Updated
	}

	projectSaved, err := database.Projects.Save(projectToSave)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to save project to database: %v", err)))
	}

	return c.JSON(http.StatusOK, projectSaved)
}
