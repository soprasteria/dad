package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ValueWithDescription represent a type that has a int value a description explaining the value
type ValueWithDescription interface {
	GetValue() int
	GetDescription() string
}

// Goal is the goal to reach for a project in deployment
type Goal struct {
	Value       int    `bson:"value" json:"value"`
	Description string `bson:"description" json:"description"`
}

// GetValue gets the value from the goal (like 10%)
func (g Goal) GetValue() int { return g.Value }

// GetDescription gets the description
func (g Goal) GetDescription() string { return g.Description }

// Progress is the progress of a project in deployment
type Progress struct {
	Value       int    `bson:"value" json:"value"`
	Description string `bson:"description" json:"description"`
}

// GetValue gets the value from the progress (like 10%)
func (p Progress) GetValue() int { return p.Value }

// GetDescription gets the description
func (p Progress) GetDescription() string { return p.Description }

// MatrixLine represents information of a depending on the functional service
type MatrixLine struct {
	Service  bson.ObjectId `bson:"service" json:"service"`
	Progress Progress      `bson:"progress" json:"progress"`
	Goal     Goal          `bson:"goal" json:"goal"`
	Comment  string        `bson:"comment" json:"comment"`
}

// Matrix represent a slice of matrix lines
type Matrix []MatrixLine

// Project represents a Sopra Steria project
type Project struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string        `bson:"name" json:"name"`
	Domain        string        `bson:"domain" json:"domain"`
	Entity        bson.ObjectId `bson:"entity" json:"entity"`
	URL           string        `bson:"url" json:"url"`
	Matrix        Matrix        `bson:"matrix" json:"matrix"`
	ServiceCenter string        `bson:"serviceCenter" json:"serviceCenter"`
	Description   string        `bson:"description" json:"description"`
	Created       time.Time     `bson:"created" json:"created"`
	Updated       time.Time     `bson:"updated" json:"updated"`
}

// UniqIDs returns the slice of Object id, where an id can appear only once
func UniqIDs(ids []bson.ObjectId) []bson.ObjectId {
	result := []bson.ObjectId{}
	seen := map[bson.ObjectId]bool{}
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			result = append(result, id)
			seen[id] = true
		}
	}
	return result
}
