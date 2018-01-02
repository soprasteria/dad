package types

import (
	"encoding/json"
	"errors"

	"github.com/soprasteria/dad/server/docktor"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FunctionalService represents the service
type FunctionalService struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string        `bson:"name" json:"name"`
	Package  string        `bson:"package" json:"package"`
	Position int           `bson:"position" json:"position"`
	Services []string      `bson:"services" json:"services"`
}

// FunctionalServiceRepo wraps all requests to database for accessing functional services
type FunctionalServiceRepo struct {
	database *mgo.Database
}

// NewFunctionalServiceRepo creates a new entites repo from database
// This FunctionalServiceRepo is wrapping all requests with database
func NewFunctionalServiceRepo(database *mgo.Database) FunctionalServiceRepo {
	return FunctionalServiceRepo{database: database}
}

func (r *FunctionalServiceRepo) col() *mgo.Collection {
	return r.database.C("functionalServices")
}

func (r *FunctionalServiceRepo) isInitialized() bool {
	return r.database != nil
}

// FindByID get the functional service by its id (string version)
func (r *FunctionalServiceRepo) FindByID(id string) (FunctionalService, error) {
	return r.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the functional service by its id (as a bson object)
func (r *FunctionalServiceRepo) FindByIDBson(id bson.ObjectId) (FunctionalService, error) {
	if !r.isInitialized() {
		return FunctionalService{}, ErrDatabaseNotInitialized
	}
	result := FunctionalService{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindAll get all functional services from the database
func (r *FunctionalServiceRepo) FindAll() ([]FunctionalService, error) {
	if !r.isInitialized() {
		return []FunctionalService{}, ErrDatabaseNotInitialized
	}
	functionalServices := []FunctionalService{}
	err := r.col().Find(bson.M{}).Sort("package", "position").All(&functionalServices)
	if err != nil {
		return []FunctionalService{}, errors.New("Can't retrieve all functional services")
	}
	return functionalServices, nil
}

// FindAllNoEmptyServices get all functional services with a service from the database
func (r *FunctionalServiceRepo) FindAllNoEmptyServices() ([]FunctionalService, error) {
	if !r.isInitialized() {
		return []FunctionalService{}, ErrDatabaseNotInitialized
	}
	functionalServices := []FunctionalService{}
	err := r.col().Find(bson.M{"$where": "this.services.length > 0"}).Sort("package", "position").All(&functionalServices)
	if err != nil {
		return []FunctionalService{}, errors.New("Can't retrieve all functional services")
	}
	return functionalServices, nil
}

// FindFunctionnalServicesDeployByProject find all deploy functional services by docktor containers
func (r *FunctionalServiceRepo) FindFunctionnalServicesDeployByProject(docktorGroup docktor.GroupDocktor) ([]FunctionalService, error) {

	functionalServices := []FunctionalService{}
	containers := []string{}
	for _, container := range docktorGroup.Containers {
		containers = append(containers, container.ServiceTitle)
	}

	jsonContainers, err := json.Marshal(containers)
	if err != nil {
		return []FunctionalService{}, errors.New("Unable to parse json docktorGroup container")
	}

	err = r.col().Find(bson.M{
		`$where`: `function () {
			if(this.services.length === 0){
				return false;
			}
			var servicesAvailable = ` + string(jsonContainers) + `
			for (var i = 0; i < this.services.length; i++) {
				if(servicesAvailable.indexOf(this.services[i]) === -1){
					return false
				}
			}
			return true
		}`,
	}).All(&functionalServices)

	return functionalServices, err
}

// Exists checks if a functional service (name and package) already exists
func (r *FunctionalServiceRepo) Exists(name, pkg string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"name":    name,
		"package": pkg,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or create the functional service in database
func (r *FunctionalServiceRepo) Save(functionalService FunctionalService) (FunctionalService, error) {
	if !r.isInitialized() {
		return FunctionalService{}, ErrDatabaseNotInitialized
	}

	if functionalService.ID.Hex() == "" {
		functionalService.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(functionalService.ID, bson.M{"$set": functionalService})
	return functionalService, err
}

// Delete the functional service
func (r *FunctionalServiceRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}
