package types

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Progress maps the progress codes to their string representation
var Progress = map[int]string{
	-1: "N/A",
	0:  "0%",
	1:  "20%",
	2:  "40%",
	3:  "60%",
	4:  "80%",
	5:  "100%",
}

// Deployed maps the progress codes to their string representation
var Deployed = map[int]string{
	-1: "No",
	0:  "Yes",
}

// Priority maps the progress codes to their string representation
var Priority = map[int]string{
	-1: "N/A",
	0:  "P0",
	1:  "P1",
	2:  "P2",
}

// MatrixLine represents information of a depending on the functional service
type MatrixLine struct {
	Service  bson.ObjectId `bson:"service" json:"service"`
	Deployed string        `bson:"deployed,omitempty" json:"deployed,omitempty"`
	Progress int           `bson:"progress" json:"progress"`
	Goal     int           `bson:"goal" json:"goal"`
	Priority string        `bson:"priority,omitempty" json:"priority,omitempty"`
	DueDate  *time.Time    `bson:"dueDate,omitempty" json:"dueDate,omitempty"`
	Comment  string        `bson:"comment" json:"comment"`
}

// Matrix represent a slice of matrix lines
type Matrix []MatrixLine

// TechnicalData contains the technical data of a project
type TechnicalData struct {
	Technologies                   []string `bson:"technologies" json:"technologies"`
	Mode                           string   `bson:"mode" json:"mode"`
	DeliverablesInVersionControl   bool     `bson:"deliverables" json:"deliverables"`
	SpecificationsInVersionControl bool     `bson:"specifications" json:"specifications"`
	SourceCodeInVersionControl     bool     `bson:"sourceCode" json:"sourceCode"`
	VersionControlSystem           string   `bson:"versionControlSystem" json:"versionControlSystem"`
	IsCDKApplicable                bool     `bson:"isCDKApplicable" json:"isCDKApplicable"`
}

// DocktorURL represents the url of the Docktor project linked to the DAD project
type DocktorURL struct {
	DocktorGroupName string `bson:"docktorGroupName" json:"docktorGroupName"`
	DocktorGroupURL  string `bson:"docktorGroupURL" json:"docktorGroupURL"`
}

// Project represents a Sopra Steria project
type Project struct {
	ID             bson.ObjectId                  `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string                         `bson:"name" json:"name"`
	Description    string                         `bson:"description" json:"description"`
	Domain         []string                       `bson:"domain" json:"domain"`
	Client         string                         `bson:"client" json:"client"`
	ProjectManager string                         `bson:"projectManager" json:"projectManager"`
	Deputies       []string                       `bson:"deputies" json:"deputies"`
	BusinessUnit   string                         `bson:"businessUnit" json:"businessUnit"`
	ServiceCenter  string                         `bson:"serviceCenter" json:"serviceCenter"`
	DocktorURL     `bson:"docktorURL" json:""`    // json is an empty string because we want to flatten the object to avoid client-side null-checks
	TechnicalData  `bson:"technicalData" json:""` // json is an empty string because we want to flatten the object to avoid client-side null-checks
	Matrix         Matrix                         `bson:"matrix" json:"matrix"`
	Created        time.Time                      `bson:"created" json:"created"`
	Updated        time.Time                      `bson:"updated" json:"updated"`
}

// Projects represents a slice of Project
type Projects []Project

// ContainsBsonID checks that a list of projects contains a certain ObjectID
func (projects Projects) ContainsBsonID(id bson.ObjectId) bool {
	for _, project := range projects {
		if project.ID == id {
			return true
		}
	}
	return false
}

// UniqIDs returns the slice of Object id, where an id can appear only once
func UniqIDs(ids []bson.ObjectId) []bson.ObjectId {
	result := []bson.ObjectId{}
	seen := map[bson.ObjectId]bool{}
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			result = append(result, id)
			seen[id] = true
		}
	}
	return result
}

// ProjectRepo wraps all requests to database for accessing entities
type ProjectRepo struct {
	database *mgo.Database
}

// NewProjectRepo creates a new projects repo from database
// This ProjectRepo is wrapping all requests with database
func NewProjectRepo(database *mgo.Database) ProjectRepo {
	return ProjectRepo{database: database}
}

func (r *ProjectRepo) col() *mgo.Collection {
	return r.database.C("projects")
}

func (r *ProjectRepo) isInitialized() bool {
	return r.database != nil
}

// FindByID get the project by its id (string version)
func (r *ProjectRepo) FindByID(id string) (Project, error) {
	return r.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the project by its id (as a bson object)
func (r *ProjectRepo) FindByIDBson(id bson.ObjectId) (Project, error) {
	if !r.isInitialized() {
		return Project{}, ErrDatabaseNotInitialized
	}
	result := Project{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindByName find a project by its name (case insensitive)
func (r *ProjectRepo) FindByName(name string) (Project, error) {
	if !r.isInitialized() {
		return Project{}, ErrDatabaseNotInitialized
	}
	result := Project{}
	regex := "^" + name + "$"
	err := r.col().Find(bson.M{"name": bson.RegEx{Pattern: regex, Options: "i"}}).One(&result)
	return result, err
}

// FindAll get all projects from the database
func (r *ProjectRepo) FindAll() ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialized
	}
	projects := []Project{}
	err := r.col().Find(bson.M{}).All(&projects)
	if err != nil {
		return []Project{}, errors.New("Can't retrieve all projects")
	}
	return projects, nil
}

func removeDuplicates(projects []Project) []Project {
	seen := map[bson.ObjectId]bool{}
	result := []Project{}

	for _, project := range projects {
		if !seen[project.ID] {
			seen[project.ID] = true
			result = append(result, project)
		}
	}
	return result
}

// FindForUser returns the projects associated to a user, handling their rights
func (r *ProjectRepo) FindForUser(user User) (Projects, error) {
	var projects []Project
	var err error

	switch user.Role {
	case AdminRole:
		projects, err = r.FindAll()
	case RIRole:
		projects, err = r.FindByEntities(user.Entities)
		if err != nil {
			return nil, err
		}

		var projectsByPMOrDeputy []Project
		projectsByPMOrDeputy, err = r.FindByProjectManagerOrDeputy(user.ID)
		if err != nil {
			return nil, err
		}

		projects = removeDuplicates(append(projects, projectsByPMOrDeputy...))
	case PMRole, DeputyRole:
		projects, err = r.FindByProjectManagerOrDeputy(user.ID)
	default:
		return nil, fmt.Errorf("Invalid role %s for user %s", user.Role, user.Username)
	}

	return projects, err
}

// FindWithDocktorGroupURL returns the projects with a no empty docktor group url
func (r *ProjectRepo) FindWithDocktorGroupURL() ([]Project, error) {
	projects := []Project{}
	err := r.col().Find(bson.M{
		"$and": []bson.M{
			bson.M{"docktorURL.docktorGroupURL": bson.M{"$exists": true, "$ne": ""}},
		},
	}).All(&projects)
	return projects, err
}

// FindModifiableForUser returns the projects associated to a user, but only projects which are modifiable by him
func (r *ProjectRepo) FindModifiableForUser(user User) (Projects, error) {
	var projects []Project
	var err error

	switch user.Role {
	case AdminRole:
		projects, err = r.FindAll()
	case RIRole:
		projects, err = r.FindByEntities(user.Entities)
		if err != nil {
			return nil, err
		}
	case PMRole, DeputyRole:
		projects, err = r.FindByProjectManagerOrDeputy(user.ID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid role %s for user %s", user.Role, user.Username)
	}

	return projects, err
}

// FindByEntities get all projects with a matching businessUnit or serviceCenter
func (r *ProjectRepo) FindByEntities(ids []bson.ObjectId) ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialized
	}

	idsString := []string{}
	for _, id := range ids {
		idsString = append(idsString, id.Hex())
	}

	projects := []Project{}
	err := r.col().Find(bson.M{
		"$or": []bson.M{
			{"businessUnit": bson.M{"$in": idsString}},
			{"serviceCenter": bson.M{"$in": idsString}},
		},
	}).All(&projects)
	if err != nil {
		return []Project{}, fmt.Errorf("Can't retrieve projects for entities %v", ids)
	}
	return projects, nil
}

// FindByProjectManagerOrDeputy get all projects with a specific project manager or deputy
func (r *ProjectRepo) FindByProjectManagerOrDeputy(id bson.ObjectId) ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialized
	}
	projectsbyPM := []Project{}
	err1 := r.col().Find(bson.M{"projectManager": id.Hex()}).All(&projectsbyPM)
	projectsbyDeputy := []Project{}
	err2 := r.col().Find(bson.M{"deputies": id.Hex()}).All(&projectsbyDeputy)
	if err1 != nil && err2 != nil {
		return []Project{}, fmt.Errorf("Can't retrieve projects for project manager and deputy %s", id)
	}
	projects := removeDuplicates(append(projectsbyPM, projectsbyDeputy...))
	return projects, nil
}

// Save updates or create the functional service in database
func (r *ProjectRepo) Save(project Project) (Project, error) {
	if !r.isInitialized() {
		return Project{}, ErrDatabaseNotInitialized
	}

	if project.ID.Hex() == "" {
		project.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(project.ID, bson.M{"$set": project})
	return project, err
}

// RemoveEntity removes an entity (businessUnit or serviceCenter) from a project
// This is used for cascade deletions
func (r *ProjectRepo) RemoveEntity(id string) error {
	if !r.isInitialized() {
		return ErrDatabaseNotInitialized
	}

	_, err := r.col().UpdateAll(
		bson.M{"businessUnit": id},
		bson.M{"$set": bson.M{"businessUnit": ""}},
	)
	if err != nil {
		return err
	}

	_, err = r.col().UpdateAll(
		bson.M{"serviceCenter": id},
		bson.M{"$set": bson.M{"serviceCenter": ""}},
	)
	return err
}

// Delete the project
func (r *ProjectRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}

// UpdateDocktorGroupURL updates Docktor Group URL to project in database
func (r *ProjectRepo) UpdateDocktorGroupURL(id bson.ObjectId, docktorGroupURL, docktorGroupName string) error {
	if !r.isInitialized() {
		return ErrDatabaseNotInitialized
	}
	return r.col().UpdateId(
		id,
		bson.M{"$set": bson.M{
			"docktorURL.docktorGroupURL":  docktorGroupURL,
			"docktorURL.docktorGroupName": docktorGroupName,
			"updated":                     time.Now(),
		}},
	)
}
