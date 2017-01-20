package types

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrganizationType identifies the type of the organization
type OrganizationType string

const (
	// BusinessUnitType is the type of organization for business units
	BusinessUnitType OrganizationType = "businessUnit"
	// ServiceCenterType is the type of organization for services centers
	ServiceCenterType OrganizationType = "serviceCenter"
)

// Organization represents an Sopra Steria organization
type Organization struct {
	ID   bson.ObjectId    `bson:"_id,omitempty" json:"id,omitempty"`
	Name string           `bson:"name" json:"name"`
	Type OrganizationType `bson:"type" json:"type"`
}

// GetID gets the ID of the organization
func (e Organization) GetID() bson.ObjectId {
	return e.ID
}

// GetOrganizationsIds get ids of a slice of organizations
func GetOrganizationsIds(organizations []Organization) []bson.ObjectId {
	ids := []bson.ObjectId{}
	for _, e := range organizations {
		ids = append(ids, e.GetID())
	}
	return ids
}

// OrganizationRepo wraps all requests to database for accessing organizations
type OrganizationRepo struct {
	database *mgo.Database
}

// NewOrganizationRepo creates a new entites repo from database
// This OrganizationRepo is wrapping all requests with database
func NewOrganizationRepo(database *mgo.Database) OrganizationRepo {
	return OrganizationRepo{database: database}
}

func (r *OrganizationRepo) col() *mgo.Collection {
	return r.database.C("organizations")
}

func (r *OrganizationRepo) isInitialized() bool {
	return r.database != nil
}

// FindByID get the organization by its id (string version)
func (r *OrganizationRepo) FindByID(id string) (Organization, error) {
	return r.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the organization by its id (as a bson object)
func (r *OrganizationRepo) FindByIDBson(id bson.ObjectId) (Organization, error) {
	if !r.isInitialized() {
		return Organization{}, ErrDatabaseNotInitialiazed
	}
	result := Organization{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindAll get all organizations from the database
func (r *OrganizationRepo) FindAll() ([]Organization, error) {
	if !r.isInitialized() {
		return []Organization{}, ErrDatabaseNotInitialiazed
	}
	organizations := []Organization{}
	err := r.col().Find(bson.M{}).All(&organizations)
	if err != nil {
		return []Organization{}, errors.New("Can't retrieve all organizations")
	}
	return organizations, nil
}

// FindAllByIDBson gets all the organizations existing with ids
func (r *OrganizationRepo) FindAllByIDBson(ids []bson.ObjectId) ([]Organization, error) {
	organizations := []Organization{}
	err := r.col().Find(bson.M{"_id": bson.M{"$in": ids}}).All(&organizations)
	if err != nil {
		return []Organization{}, errors.New("Can't retrieve all organizations")
	}
	return organizations, nil
}

// Exists checks if an organization (name) already exists
func (r *OrganizationRepo) Exists(name string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"name": name,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or create the organization in database
func (r *OrganizationRepo) Save(organization Organization) (Organization, error) {
	if !r.isInitialized() {
		return Organization{}, ErrDatabaseNotInitialiazed
	}

	if organization.ID.Hex() == "" {
		organization.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(organization.ID, bson.M{"$set": organization})
	return organization, err
}

// Delete the organization
func (r *OrganizationRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}
