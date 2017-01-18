package types

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// IsDatabase is an interface representing a database accessing mongod documents
type IsDatabase interface {
	col() *mgo.Collection
	isInitialized() bool
}

// IsDocument is an interface representing a data being a collection in Mongo
type IsDocument interface {
	GetID() bson.ObjectId
}

// BasicDelete is a basic delete of a collection document
func BasicDelete(collection IsDatabase, id bson.ObjectId) (bson.ObjectId, error) {
	if !collection.isInitialized() {
		return bson.ObjectIdHex(""), ErrDatabaseNotInitialiazed
	}

	err := collection.col().RemoveId(id)
	return id, err
}
