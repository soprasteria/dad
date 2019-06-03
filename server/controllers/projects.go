package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/matcornic/hermes"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/soprasteria/dad/server/docktor"
	"github.com/soprasteria/dad/server/email"
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

// sendEmail sets the email body and send it
func sendEmail(project types.Project, to types.User, userRepo types.UserRepo, name, address string) {

	var option = email.SendOptions{
		To: []mail.Address{
			{
				Name:    to.DisplayName,
				Address: to.Email,
			},
		},

		ToCc: []mail.Address{
			{
				Name:    name,
				Address: address,
			},
		},

		Subject: "D.A.D - Project deleted",

		Body: hermes.Email{
			Body: hermes.Body{
				Intros: []string{
					"WARNING: This project was linked with a Docktor."},

				Title: "Project " + project.Name + " deleted!",

				Dictionary: []hermes.Entry{
					{Key: "URL Docktor", Value: project.DocktorURL.DocktorGroupURL},
				},
			},
		}}

	// Add the RI as email receiver
	entityIDs := []bson.ObjectId{bson.ObjectIdHex(project.BusinessUnit)}
	for _, sC := range project.ServiceCenter {
		entityIDs[len(entityIDs)] = bson.ObjectIdHex(sC)
	}
	riUsers, err := userRepo.FindRIWithEntity(entityIDs)
	if err != nil {
		log.Error("Error while retrieving the users whose matching with serviceCenter/businessUnit IDs or with RI's role: ", err)
	} else {
		for _, u := range riUsers {
			option.To = append(option.To, mail.Address{Name: u.DisplayName, Address: u.Email})
		}
	}

	// Add the PM as email receiver
	if project.ProjectManager != "" {
		projectManager, err := userRepo.FindByID(project.ProjectManager)
		if err != nil {
			log.Error("Error while retrieving the project manager stats from the UserRepo: ", err)
		} else {
			option.To = append(option.To, mail.Address{Name: projectManager.DisplayName, Address: projectManager.Email})
		}
	}

	// Add the deputies as email receiver
	for _, d := range project.Deputies {
		deputies, err := userRepo.FindByID(d)
		if err != nil {
			log.Error("Error while retrieving the deputies from the UserRepo: ", err)
		} else {
			option.To = append(option.To, mail.Address{Name: deputies.DisplayName, Address: deputies.Email})
		}
	}

	// Send the email
	errorSend := email.Send(option)

	if errorSend != nil {
		log.Error("Error while sending an email", errorSend)
	}
}

// Delete project from database
func (p *Projects) Delete(c echo.Context) error {
	database := c.Get("database").(*mongo.DadMongo)
	id := c.Param("id")
	authUser := c.Get("authuser").(types.User)
	userRepo := database.Users

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

	// Get the project's stats before deleting it
	projectStats, err := database.Projects.FindByID(id)
	if err != nil {
		log.Error("Error while retrieving the project from database", err)
	}

	// Deleting the project
	res, err := database.Projects.Delete(bson.ObjectIdHex(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewErr(fmt.Sprintf("Error while removing project: %v", err)))
	}

	// checks if deleted project had a linked Docktor URL.
	if projectStats.DocktorURL.DocktorGroupURL != "" {
		sendEmail(projectStats, authUser, userRepo, viper.GetString("name.receiver"), viper.GetString("admin.email"))
	}

	return c.JSON(http.StatusOK, res)
}

func canAddEntityToProject(entityToSet string, entityFromDB []string, authUser types.User) bool {
	if authUser.Role == types.AdminRole {
		return true
	}

	if !bson.IsObjectIdHex(entityToSet) {
		return false
	}

	// If the user is a RI, he can only add an entity if:
	// * it's the currently assigned entity of the project
	for _, eDB := range entityFromDB {
		if bson.IsObjectIdHex(eDB) && bson.ObjectIdHex(entityToSet) == bson.ObjectIdHex(eDB) {
			return true
		}
	}
	// * it's one of his own assigned entities
	for _, allowedEntity := range authUser.Entities {
		if bson.ObjectIdHex(entityToSet) == allowedEntity {
			return true
		}
	}
	return false
}

func validateEntity(entityRepo types.EntityRepo, entityToSet, entityFromDB []string, entityType types.EntityType, authUser types.User) (int, string) {
	if len(entityToSet) > 0 {
		for _, eS := range entityToSet {
			// Retrieve entity check if exist
			entity, err := entityRepo.FindByID(eS)
			if err == mgo.ErrNotFound {
				return http.StatusBadRequest, fmt.Sprintf("The %s %s does not exist", entityType, eS)
			} else if err != nil {
				return http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve %s %s from database: %v", entityType, eS, err)
			}
			// Check type if BU or service center ...
			if entity.Type != entityType {
				return http.StatusBadRequest, fmt.Sprintf("The entity %s (%s) is not of type %s but  %s", entity.Name, eS, entityType, entity.Type)
			}

			// Check if you have the access to change the entity
			if !canAddEntityToProject(eS, entityFromDB, authUser) {
				return http.StatusBadRequest, fmt.Sprintf("You can't add the entity %s to a project", eS)
			}
		}
	}
	return http.StatusOK, ""
}

func validateEntities(entityRepo types.EntityRepo, projectToSave, projectFromDB types.Project, authUser types.User) (int, string) {
	if projectToSave.BusinessUnit == "" && len(projectToSave.ServiceCenter) == 0 {
		return http.StatusBadRequest, "At least one of the business unit and service center fields is mandatory"
	}

	// Check if BusinessUnit is set because if projectToSave.BusinessUnit is nil, len([]string{projectToSave.BusinessUnit}) will return 1
	if projectToSave.BusinessUnit != "" {
		// If a business unit is provided, check it exists in the entity collection
		if statusCode, errMessage := validateEntity(entityRepo, []string{projectToSave.BusinessUnit}, []string{projectFromDB.BusinessUnit}, types.BusinessUnitType, authUser); errMessage != "" {
			return statusCode, errMessage
		}
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
	// Get the project ID from the URL (used to distinguish between create and update)
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
			}).Warn("User isn't allowed to view the project")
			return c.JSON(http.StatusForbidden, types.NewErr(fmt.Sprintf("User %s isn't allowed to update the project", authUser.Username)))
		}

		projectFromDB, err = database.Projects.FindByID(id)
		if err != nil || projectFromDB.ID.Hex() == "" {
			return c.JSON(http.StatusBadRequest, types.NewErr("Trying to modify a non existing project"))
		}
	}

	var projectToSave types.Project
	// Get project from the body
	err := c.Bind(&projectToSave)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.NewErr(fmt.Sprintf("Posted project is not valid: %v", err)))
	}

	saveProjectData, httpStatus, err := p.createProjectToSave(database, id, projectToSave, authUser, projectFromDB)
	if err != nil {
		return c.JSON(httpStatus, types.NewErr(fmt.Sprintf("Error while creating the new project to save: %v", err.Error())))
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
				return
			}
			defer database.Session.Close()

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

func (p *Projects) createProjectToSave(database *mongo.DadMongo, id string, projectToSave types.Project, authUser types.User, projectFromDB types.Project) (SaveProjectData, int, error) {
	log.WithField("project", projectToSave).Info("Received project to save")

	if projectToSave.Name == "" {
		return SaveProjectData{}, http.StatusBadRequest, errors.New("The name field cannot be empty")
	}

	// Not possible to create or update a project with a name already used by another one project
	existingProject, err := database.Projects.FindByName(projectToSave.Name)
	if err != nil {
		if err != mgo.ErrNotFound {
			return SaveProjectData{}, http.StatusInternalServerError, fmt.Errorf("Can't check whether the project exist in database: %v", err)
		}
	} else if existingProject.ID != projectToSave.ID {
		return SaveProjectData{}, http.StatusBadRequest, fmt.Errorf("Another project already exists with the same name %q", existingProject.Name)
	}

	// check if id is valid (for project creation)
	if projectToSave.ID.Valid() {
		existingProject, err = database.Projects.FindByIDBson(projectToSave.ID)
		if err != nil {
			if err != mgo.ErrNotFound {
				return SaveProjectData{}, http.StatusInternalServerError, fmt.Errorf("Can't check whether the project exist in database: %v", err)
			}
		}
	}

	modifiedDetails := projectToSave.Name != existingProject.Name ||
		strings.Join(projectToSave.Domain, ";") != strings.Join(existingProject.Domain, ";") ||
		projectToSave.ProjectManager != existingProject.ProjectManager ||
		!reflect.DeepEqual(projectToSave.ServiceCenter, existingProject.ServiceCenter) ||
		projectToSave.BusinessUnit != existingProject.BusinessUnit ||
		projectToSave.DocktorGroupURL != existingProject.DocktorGroupURL

	// A Project Manager or Deputy can't update details, if any of the details has changed it's an issue and we shouldn't update the project
	if (authUser.Role == types.PMRole || authUser.Role == types.DeputyRole) && modifiedDetails {

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
		return SaveProjectData{}, http.StatusBadRequest, errors.New("Project managers and deputies are not allowed to update project details")
	} else if authUser.Role == types.RIRole {
		if projectToSave.Mode != existingProject.Mode {
			return SaveProjectData{}, http.StatusBadRequest, errors.New("RIs are not allowed to update deployment mode")
		}
	}

	// Check rights to add entities to the project
	httpStatusCode, errorMessage := validateEntities(database.Entities, projectToSave, projectFromDB, authUser)
	if errorMessage != "" {
		return SaveProjectData{}, httpStatusCode, errors.New(errorMessage)
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
	return NewSaveProjectData(existingProject, projectToSave), 0, nil
}

// updateDocktorGroupName updates the Docktor Group Name in saved project
// It gets the Group name from Docktor Group URL by fetching Docktor API directly
func (p *Projects) updateDocktorGroupName(database *mongo.DadMongo, idProject bson.ObjectId, docktorGroupURL string) error {

	// Call Docktor API to get the real name of the group
	docktorAPI, err := docktor.NewExternalAPI(
		viper.GetString("docktor.addr"),
		viper.GetString("docktor.user"),
		viper.GetString("docktor.password"),
		viper.GetBool("docktor.ldap"),
	)
	if err != nil {
		return err
	}
	// Parse Docktor URL to get the Docktor group ID
	idDocktorGroup, err := docktorAPI.GetGroupIDFromURL(docktorGroupURL)
	if err != nil {
		return err
	}
	group, err := docktorAPI.GetGroup(idDocktorGroup)
	if err != nil {
		return err
	}
	// Update project in database
	err = database.Projects.UpdateDocktorGroupURL(idProject, docktorGroupURL, group.Name)
	if err != nil {
		return fmt.Errorf("Unable to update project in Mongo database because: %v", err.Error())
	}

	return nil
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
