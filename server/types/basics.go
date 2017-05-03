package types

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrorMsg is a json formated error
type ErrorMsg struct {
	Message string `json:"message"`
}

// IsDatabase is an interface representing a database accessing mongod documents
type IsDatabase interface {
	col() *mgo.Collection
	isInitialized() bool
}

// IsDatabaseWithIndexes is an interface representing a database which needs indexes to be created
type IsDatabaseWithIndexes interface {
	IsDatabase
	CreateIndexes() error
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

// NewErr is a function used to format errors into json
func NewErr(message string) ErrorMsg {
	return ErrorMsg{Message: message}
}
