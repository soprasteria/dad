package types

import "gopkg.in/mgo.v2/bson"

// Entity represents an Sopra Steria entity
type Entity struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name string        `bson:"name" json:"name"`
}
