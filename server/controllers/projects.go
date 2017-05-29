package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/docktor"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
)

// Projects is the controller type
type Projects struct {
}

// SaveProjectData is the Struct used in order to manage project to Save easily
type SaveProjectData struct {
	existingProject types.Project
	projectToSave   types.Project
}

// NewSaveProjectData is SaveProjectData Constructor
func NewSaveProjectData(existingProject types.Project, projectToSave types.Project) SaveProjectData {
	return SaveProjectData{existingProject: existingProject, projectToSave: projectToSave}
}

// GetAll functional services from database
func (p *Projects) GetAll(c echo.Context) error {
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

// Get project from database. The project was stored by a middleware which use id to get project informations
func (p *Projects) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("project"))
}

// GetIndicators corresponding to a specific project in database. The project was stored by a middleware which use id to get project informations
func (p *Projects) GetIndicators(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	project := c.Get("project").(types.Project)
	indicators, err := database.UsageIndicators.FindAllFromGroup(project.DocktorGroupName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while retrieving the indicators of the project %s: %v", project.DocktorGroupName, err.Error())))
	}
	return c.JSON(http.StatusOK, indicators)
}

// Delete project from database
func (p *Projects) Delete(c echo.Context) error {
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
func (p *Projects) Save(c echo.Context) error {
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

	saveProjectData, err := p.createProjectToSave(database, c, authUser, projectFromDB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while creating the new project to save: %v", err)))
	}

	projectSaved, err := database.Projects.Save(saveProjectData.projectToSave)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to save project to database: %v", err)))
	}
	if projectSaved.DocktorGroupURL != "" && saveProjectData.existingProject.DocktorGroupURL != projectSaved.DocktorGroupURL {
		// Updates Docktor group name from url in background, when url changed
		// because we don't want to block the project update with calls to external APIs.
		go func() {
			logFields := log.Fields{
				"dad.project.id":                       projectSaved.ID,
				"dad.project.name":                     projectSaved.Name,
				"dad.project.docktorGroupURL":          projectSaved.DocktorGroupURL,
				"dad.project.previousDocktorGroupName": projectSaved.DocktorGroupName,
			}
			log.WithFields(logFields).Debug("Updating DocktorGroupURL and Name to DAD project...")

			// Open new Mongo session because function is called in a goroutine
			database, err := mongo.Get()
			if err != nil {
				log.WithField("database", database).WithError(err).Error("Unable to open a connection to the database")
			}

			err = p.updateDocktorGroupName(database, projectSaved.ID, projectSaved.DocktorGroupURL)
			if err != nil {
				log.WithFields(logFields).WithError(err).Error("Unable to fetch and/or save Docktor Group Name to the project")
			} else {
				log.WithFields(logFields).Debug("Saved DocktorGroupURL and Name to DAD project")
			}
		}()
	}

	log.WithFields(log.Fields{
		"id":   projectSaved.ID,
		"name": projectSaved.Name,
	}).Debug("Project is saved")

	return c.JSON(http.StatusOK, projectSaved)
}

func (p *Projects) createProjectToSave(database *mongo.DadMongo, c echo.Context, authUser types.User, projectFromDB types.Project) (SaveProjectData, error) {
	var projectToSave types.Project
	// Get project from body
	err := c.Bind(&projectToSave)
	if err != nil {
		return SaveProjectData{}, c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted project is not valid: %v", err)))
	}

	log.WithField("project", projectToSave).Info("Received project to save")

	if projectToSave.Name == "" {
		return SaveProjectData{}, c.JSON(http.StatusBadRequest, types.NewErr("The name field cannot be empty"))
	}

	// Not possible to create or update a project with a name already used by another one project
	existingProject, err := database.Projects.FindByName(projectToSave.Name)
	if err != nil {
		if err != mgo.ErrNotFound {
			return SaveProjectData{}, c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Can't check whether the project exist in database: %v", err)))
		}
	} else if existingProject.ID != projectToSave.ID {
		return SaveProjectData{}, c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Another project already exists with the same name %q", existingProject.Name)))
	}

	modifiedDetails := projectToSave.Name != existingProject.Name ||
		strings.Join(projectToSave.Domain, ";") != strings.Join(existingProject.Domain, ";") ||
		projectToSave.ProjectManager != existingProject.ProjectManager ||
		projectToSave.ServiceCenter != existingProject.ServiceCenter ||
		projectToSave.BusinessUnit != existingProject.BusinessUnit ||
		projectToSave.DocktorGroupURL != existingProject.DocktorGroupURL

	var id = c.Param("id")
	// A Project Manager can't update details, if any of the details has changed it's an issue and we shouldn't update the project
	if authUser.Role == types.CPRole && modifiedDetails {
		log.WithFields(log.Fields{
			"username":                        authUser.Username,
			"role":                            authUser.Role,
			"projectID":                       id,
			"projectToSave.Name":              projectToSave.Name,
			"existingProject.Name":            existingProject.Name,
			"projectToSave.Domain":            projectToSave.Domain,
			"existingProject.Domain":          existingProject.Domain,
			"projectToSave.ProjectManager":    projectToSave.ProjectManager,
			"existingProject.ProjectManager":  existingProject.ProjectManager,
			"projectToSave.ServiceCenter":     projectToSave.ServiceCenter,
			"existingProject.ServiceCenter":   existingProject.ServiceCenter,
			"projectToSave.BusinessUnit":      projectToSave.BusinessUnit,
			"existingProject.BusinessUnit":    existingProject.BusinessUnit,
			"projectToSave.DocktorGroupURL":   projectToSave.DocktorGroupURL,
			"existingProject.DocktorGroupURL": existingProject.DocktorGroupURL,
		}).Warn("User isn't allowed to update the project")
		return SaveProjectData{}, c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("A project manager isn't allowed to update project details")))
	}

	// Check rights to add entities to the project
	httpStatusCode, errorMessage := validateEntities(database.Entities, projectToSave, projectFromDB, authUser)
	if errorMessage != "" {
		return SaveProjectData{}, c.JSON(httpStatusCode, types.NewErr(errorMessage))
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

	// We empty the docktor group name because the URL has changed. The group name will be automatically recalculated using the new URL
	if existingProject.DocktorGroupURL != projectToSave.DocktorGroupURL {
		projectToSave.DocktorGroupName = ""
	}
	return NewSaveProjectData(existingProject, projectToSave), nil
}

// updateDocktorGroupName updates the Docktor Group Name in saved project
// It gets the Group name from Docktor Group URL by fetching Docktor API directly
func (p *Projects) updateDocktorGroupName(database *mongo.DadMongo, idProject bson.ObjectId, docktorGroupURL string) error {

	// Parse Docktor URL to get the Docktor group ID
	idDocktorGroup, err := getGroupIDFromURL(docktorGroupURL)
	if err != nil {
		return err
	}
	// Call Docktor API to get the real name of the group
	docktorAPI, err := docktor.NewExternalAPI(
		viper.GetString("docktor.addr"),
		viper.GetString("docktor.user"),
		viper.GetString("docktor.password"),
	)
	if err != nil {
		return err
	}
	group, err := docktorAPI.GetGroup(idDocktorGroup)
	if err != nil {
		return err
	}
	// Update project in database
	err = database.Projects.UpdateDocktorGroupURL(idProject, docktorGroupURL, group.Title)
	if err != nil {
		return fmt.Errorf("Unable to update project in Mongo database because: %v", err.Error())
	}

	return nil
}

// getGroupIDFromURL returns the Docktor group ID from its URL
// URL is expected to be format : http://<docktor-host>/groups/<id>
func getGroupIDFromURL(docktorURL string) (string, error) {
	u, err := url.ParseRequestURI(docktorURL)
	if err != nil {
		return "", fmt.Errorf("docktorGroupURL is not a valid URL. Expected 'http://<docktor>/groups/<id>', Got '%v'", docktorURL)
	}
	path := strings.Split(u.Path, "/")
	if len(path) == 0 {
		return "", fmt.Errorf("Unable to get project id from URL. Expected 'http://<docktor>/groups/<id>', Got '%v'", u.Path)
	}
	id := path[len(path)-1]
	if id == "" {
		return "", fmt.Errorf("Unable to get project id from URL parsed path : %v. URL=%v", path, u.Path)
	}
	return id, nil
}

// UpdateDocktorInfo updates docktor info of a specific project
func (p *Projects) UpdateDocktorInfo(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")

	// Get project from body
	var projectToSave types.Project

	err := c.Bind(&projectToSave)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted project is not valid: %v", err)))
	}

	log.WithFields(log.Fields{
		"id":              id,
		"docktorGroupURL": projectToSave.DocktorGroupURL,
	}).Debug("Updating Docktor Group for given project...")

	err = p.updateDocktorGroupName(database, bson.ObjectIdHex(id), projectToSave.DocktorGroupURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to update docktor info to database: %v", err)))
	}

	project, err := database.Projects.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Failed to get the updated project to database: %v", err)))
	}

	log.WithFields(log.Fields{
		"id":                       id,
		"project.name":             project.Name,
		"project.docktorGroupURL":  project.DocktorGroupURL,
		"project.docktorGroupName": project.DocktorGroupName,
	}).Debug("Updated Docktor Group for given project")

	return c.JSON(http.StatusOK, project)
}
