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

// MatrixLine represents information of a depending on the functional service
type MatrixLine struct {
	Service  bson.ObjectId `bson:"service" json:"service"`
	Progress int           `bson:"progress" json:"progress"`
	Goal     int           `bson:"goal" json:"goal"`
	Comment  string        `bson:"comment" json:"comment"`
}

// Matrix represent a slice of matrix lines
type Matrix []MatrixLine

// Project represents a Sopra Steria project
type Project struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string        `bson:"name" json:"name"`
	Domain         string        `bson:"domain" json:"domain"`
	ProjectManager bson.ObjectId `bson:"projectManager,omitempty" json:"projectManager,omitempty"`
	BusinessUnit   bson.ObjectId `bson:"businessUnit,omitempty" json:"businessUnit,omitempty"`
	ServiceCenter  bson.ObjectId `bson:"serviceCenter,omitempty" json:"serviceCenter,omitempty"`
	URL            string        `bson:"url" json:"url"`
	Matrix         Matrix        `bson:"matrix" json:"matrix"`
	Description    string        `bson:"description" json:"description"`
	Created        time.Time     `bson:"created" json:"created"`
	Updated        time.Time     `bson:"updated" json:"updated"`
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
		return Project{}, ErrDatabaseNotInitialiazed
	}
	result := Project{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindAll get all projects from the database
func (r *ProjectRepo) FindAll() ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialiazed
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

		projectsByPM, err := r.FindByProjectManager(user.ID)
		if err != nil {
			return nil, err
		}

		projects = removeDuplicates(append(projects, projectsByPM...))
	case CPRole:
		projects, err = r.FindByProjectManager(user.ID)
	default:
		return nil, fmt.Errorf("Invalid role %s for user %s", user.Role, user.Username)
	}

	return projects, err
}

// FindByEntities get all projects with a matching businessUnit or serviceCenter
func (r *ProjectRepo) FindByEntities(ids []bson.ObjectId) ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialiazed
	}
	projects := []Project{}
	err := r.col().Find(bson.M{
		"$or": bson.M{
			"businessUnit":  bson.M{"$in": ids},
			"serviceCenter": bson.M{"$in": ids},
		},
	}).All(&projects)
	if err != nil {
		return []Project{}, fmt.Errorf("Can't retrieve projects for entities %v", ids)
	}
	return projects, nil
}

// FindByProjectManager get all projects with a specific project manager
func (r *ProjectRepo) FindByProjectManager(id bson.ObjectId) ([]Project, error) {
	if !r.isInitialized() {
		return []Project{}, ErrDatabaseNotInitialiazed
	}
	projects := []Project{}
	err := r.col().Find(bson.M{"projectManager": id}).All(&projects)
	if err != nil {
		return []Project{}, fmt.Errorf("Can't retrieve projects for project manager %s", id)
	}
	return projects, nil
}

// Save updates or create the functionnal service in database
func (r *ProjectRepo) Save(project Project) (Project, error) {
	if !r.isInitialized() {
		return Project{}, ErrDatabaseNotInitialiazed
	}

	if project.ID.Hex() == "" {
		project.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(project.ID, bson.M{"$set": project})
	return project, err
}

// Delete the project
func (r *ProjectRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}
