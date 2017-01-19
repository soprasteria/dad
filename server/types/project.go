package types

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RichValue is a value with a description
type RichValue struct {
	Value       int    `bson:"value" json:"value"`
	Description string `bson:"description" json:"description"`
}

// MatrixLine represents information of a depending on the functional service
type MatrixLine struct {
	Service  bson.ObjectId `bson:"service" json:"service"`
	Progress RichValue     `bson:"progress" json:"progress"`
	Goal     RichValue     `bson:"goal" json:"goal"`
	Comment  string        `bson:"comment" json:"comment"`
}

// Matrix represent a slice of matrix lines
type Matrix []MatrixLine

// Project represents a Sopra Steria project
type Project struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string        `bson:"name" json:"name"`
	Domain        string        `bson:"domain" json:"domain"`
	Entity        bson.ObjectId `bson:"entity,omitempty" json:"entity,omitempty"`
	ServiceCenter bson.ObjectId `bson:"serviceCenter,omitempty" json:"serviceCenter,omitempty"`
	URL           string        `bson:"url" json:"url"`
	Matrix        Matrix        `bson:"matrix" json:"matrix"`
	Description   string        `bson:"description" json:"description"`
	Created       time.Time     `bson:"created" json:"created"`
	Updated       time.Time     `bson:"updated" json:"updated"`
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

// FindAllByIDBson gets all the projects existing with ids
func (r *ProjectRepo) FindAllByIDBson(ids []bson.ObjectId) ([]Project, error) {
	projects := []Project{}
	err := r.col().Find(bson.M{"_id": bson.M{"$in": ids}}).All(&projects)
	if err != nil {
		return []Project{}, errors.New("Can't retrieve all projects")
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
