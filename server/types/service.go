package types

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// FunctionalService represents the service
type FunctionalService struct {
	ID                    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name                  string        `bson:"name" json:"name"`
	Translations          Translations  `bson:"translations" json:"translations"`
	Package               string        `bson:"package" json:"package"`
	Position              int           `bson:"position" json:"position"`
	Services              []string      `bson:"services" json:"services"`
	DeclarativeDeployment bool          `bson:"declarativeDeployement" json:"declarativeDeployement"`
}

// isAssociatedWithAtLeastGivenService return true when at least one service in given parameter is found in the functional service.
// Search is executed by ignoring case.
func (fs FunctionalService) isAssociatedWithAtLeastGivenService(serviceNames []string) bool {
	for _, service := range fs.Services {
		for _, name := range serviceNames {
			if strings.ToLower(name) == strings.ToLower(service) {
				return true
			}
		}
	}
	return false
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

// FindFunctionalServicesDeployByServices find all functional services associated to
func (r *FunctionalServiceRepo) FindFunctionalServicesDeployByServices(services []string) ([]FunctionalService, error) {

	functionalServices := []FunctionalService{}

	allFunctionalServices, err := r.FindAll()
	if err != nil {
		return nil, errors.New("Unable to get all functional services from database")
	}

	for _, s := range allFunctionalServices {
		if s.isAssociatedWithAtLeastGivenService(services) {
			functionalServices = append(functionalServices, s)
		}
	}

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
