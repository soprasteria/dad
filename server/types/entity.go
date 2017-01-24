package types

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// EntityType identifies the type of the entity
type EntityType string

const (
	// BusinessUnitType is the type of entity for business units
	BusinessUnitType EntityType = "businessUnit"
	// ServiceCenterType is the type of entity for services centers
	ServiceCenterType EntityType = "serviceCenter"
)

// Entity represents an Sopra Steria entity
type Entity struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name string        `bson:"name" json:"name"`
	Type EntityType    `bson:"type" json:"type"`
}

// GetID gets the ID of the entity
func (e Entity) GetID() bson.ObjectId {
	return e.ID
}

// GetEntitiesIds get ids of a slice of entities
func GetEntitiesIds(entities []Entity) []bson.ObjectId {
	ids := []bson.ObjectId{}
	for _, e := range entities {
		ids = append(ids, e.GetID())
	}
	return ids
}

// EntityRepo wraps all requests to database for accessing entities
type EntityRepo struct {
	database *mgo.Database
}

// NewEntityRepo creates a new entites repo from database
// This EntityRepo is wrapping all requests with database
func NewEntityRepo(database *mgo.Database) EntityRepo {
	return EntityRepo{database: database}
}

func (r *EntityRepo) col() *mgo.Collection {
	return r.database.C("entities")
}

func (r *EntityRepo) isInitialized() bool {
	return r.database != nil
}

// FindByID get the entity by its id (string version)
func (r *EntityRepo) FindByID(id string) (Entity, error) {
	if !bson.IsObjectIdHex(id) {
		return Entity{}, ErrInvalidEntityID
	}
	return r.FindByIDBson(bson.ObjectIdHex(id))
}

// FindByIDBson get the entity by its id (as a bson object)
func (r *EntityRepo) FindByIDBson(id bson.ObjectId) (Entity, error) {
	if !r.isInitialized() {
		return Entity{}, ErrDatabaseNotInitialiazed
	}
	result := Entity{}
	err := r.col().FindId(id).One(&result)
	return result, err
}

// FindAll get all entities from the database
func (r *EntityRepo) FindAll() ([]Entity, error) {
	if !r.isInitialized() {
		return []Entity{}, ErrDatabaseNotInitialiazed
	}
	entities := []Entity{}
	err := r.col().Find(bson.M{}).All(&entities)
	if err != nil {
		return []Entity{}, errors.New("Can't retrieve all entities")
	}
	return entities, nil
}

// FindAllByIDBson gets all the entities existing with ids
func (r *EntityRepo) FindAllByIDBson(ids []bson.ObjectId) ([]Entity, error) {
	entities := []Entity{}
	err := r.col().Find(bson.M{"_id": bson.M{"$in": ids}}).All(&entities)
	if err != nil {
		return []Entity{}, errors.New("Can't retrieve all entities")
	}
	return entities, nil
}

// Exists checks if an entity (name) already exists
func (r *EntityRepo) Exists(name string) (bool, error) {
	nb, err := r.col().Find(bson.M{
		"name": name,
	}).Count()

	if err != nil {
		return true, err
	}
	return nb != 0, nil
}

// Save updates or create the entity in database
func (r *EntityRepo) Save(entity Entity) (Entity, error) {
	if !r.isInitialized() {
		return Entity{}, ErrDatabaseNotInitialiazed
	}

	if entity.ID.Hex() == "" {
		entity.ID = bson.NewObjectId()
	}

	_, err := r.col().UpsertId(entity.ID, bson.M{"$set": entity})
	return entity, err
}

// Delete the entity
func (r *EntityRepo) Delete(id bson.ObjectId) (bson.ObjectId, error) {
	return BasicDelete(r, id)
}
