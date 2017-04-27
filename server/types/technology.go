package types

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Technology represents a technology (ie. programming language, software suite, etc.)
type Technology struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name string        `bson:"name" json:"name"`
}

// TechnologyRepo wraps all requests to database for accessing technologies
type TechnologyRepo struct {
	database *mgo.Database
}

// NewTechnologyRepo creates a new technologies repo from database
// This TechnologyRepo is wrapping all requests with database
func NewTechnologyRepo(database *mgo.Database) TechnologyRepo {
	return TechnologyRepo{database: database}
}

func (r *TechnologyRepo) col() *mgo.Collection {
	return r.database.C("technologies")
}

func (r *TechnologyRepo) isInitialized() bool {
	return r.database != nil
}

// FindAll get all technologies from the database
func (r *TechnologyRepo) FindAll() ([]Technology, error) {
	if !r.isInitialized() {
		return []Technology{}, ErrDatabaseNotInitialiazed
	}
	technologies := []Technology{}
	err := r.col().Find(bson.M{}).All(&technologies)
	if err != nil {
		return []Technology{}, errors.New("Can't retrieve all technologies")
	}
	return technologies, nil
}

// Exists checks if a technology (name) already exists
func (r *TechnologyRepo) Exists(name string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"name": name,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or creates the technology in database
func (r *TechnologyRepo) Save(technology Technology) (Technology, error) {
	if !r.isInitialized() {
		return Technology{}, ErrDatabaseNotInitialiazed
	}

	if technology.ID.Hex() == "" {
		technology.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(technology.ID, bson.M{"$set": technology})
	return technology, err
}
