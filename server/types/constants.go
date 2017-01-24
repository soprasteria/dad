package types

import "errors"

var (
	// ErrDatabaseNotInitialiazed occurs when the database is not initialized and no request can be executed
	ErrDatabaseNotInitialiazed = errors.New("Database is not initialized")
	// ErrInvalidUserID occurs when the user id is not a valid objectID Hex
	ErrInvalidUserID = errors.New("Invalid User ID")
	// ErrInvalidEntityID occurs when the entity id is not a valid objectID Hex
	ErrInvalidEntityID = errors.New("Invalid Entity ID")
)
