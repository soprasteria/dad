package types

import "errors"

var (
	// ErrDatabaseNotInitialiazed occurs when the database is not initialized and no request can be executed
	ErrDatabaseNotInitialiazed = errors.New("Database is not initialized")
)
