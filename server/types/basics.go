package types

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrorMsg is a json formated error
type ErrorMsg struct {
	Message string `json:"message"`
}

// IsCollection is an interface representing a collection accessing mongod documents
type IsCollection interface {
	col() *mgo.Collection
	isInitialized() bool
}

// IsCollectionWithIndexes is an interface representing a collection which needs indexes to be created
type IsCollectionWithIndexes interface {
	IsCollection
	CreateIndexes() error
}

// IsDocument is an interface representing a data being in a collection in Mongo
type IsDocument interface {
	GetID() bson.ObjectId
}

// BasicDelete is a basic delete of a collection document
func BasicDelete(collection IsCollection, id bson.ObjectId) (bson.ObjectId, error) {
	if !collection.isInitialized() {
		return bson.ObjectIdHex(""), ErrDatabaseNotInitialized
	}

	err := collection.col().RemoveId(id)
	return id, err
}

// NewErr is a function used to format errors into json
func NewErr(message string) ErrorMsg {
	return ErrorMsg{Message: message}
}
