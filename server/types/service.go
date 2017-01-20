package types

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FunctionnalService represents the service
type FunctionnalService struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string        `bson:"name" json:"name"`
	Package string        `bson:"package" json:"package"`
}

// FunctionnalServiceRepo wraps all requests to database for accessing functionnal services
type FunctionnalServiceRepo struct {
	database *mgo.Database
}

// NewFunctionnalServiceRepo creates a new entites repo from database
// This FunctionnalServiceRepo is wrapping all requests with database
func NewFunctionnalServiceRepo(database *mgo.Database) FunctionnalServiceRepo {
	return FunctionnalServiceRepo{database: database}
}

func (r *FunctionnalServiceRepo) col() *mgo.Collection {
	return r.database.C("functionnalServices")
}

func (r *FunctionnalServiceRepo) isInitialized() bool {
	return r.database != nil
}

// FindByID get the functionnal service by its id (string version)
func (r *FunctionnalServiceRepo) FindByID(id string) (FunctionnalService, error) {
	return r.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the functionnal service by its id (as a bson object)
func (r *FunctionnalServiceRepo) FindByIDBson(id bson.ObjectId) (FunctionnalService, error) {
	if !r.isInitialized() {
		return FunctionnalService{}, ErrDatabaseNotInitialiazed
	}
	result := FunctionnalService{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindAll get all functionnal services from the database
func (r *FunctionnalServiceRepo) FindAll() ([]FunctionnalService, error) {
	if !r.isInitialized() {
		return []FunctionnalService{}, ErrDatabaseNotInitialiazed
	}
	functionnalServices := []FunctionnalService{}
	err := r.col().Find(bson.M{}).All(&functionnalServices)
	if err != nil {
		return []FunctionnalService{}, errors.New("Can't retrieve all functionnal services")
	}
	return functionnalServices, nil
}

// Exists checks if a functionnal service (name and package) already exists
func (r *FunctionnalServiceRepo) Exists(name, pkg string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"name":    name,
		"package": pkg,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or create the functionnal service in database
func (r *FunctionnalServiceRepo) Save(functionnalService FunctionnalService) (FunctionnalService, error) {
	if !r.isInitialized() {
		return FunctionnalService{}, ErrDatabaseNotInitialiazed
	}

	if functionnalService.ID.Hex() == "" {
		functionnalService.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(functionnalService.ID, bson.M{"$set": functionnalService})
	return functionnalService, err
}

// Delete the functionnal service
func (r *FunctionnalServiceRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}
