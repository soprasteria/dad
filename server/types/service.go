package types

import "gopkg.in/mgo.v2/bson"

// FunctionnalService represents the service
type FunctionnalService struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string        `bson:"name" json:"name"`
	Package string        `bson:"package" json:"package"`
}
