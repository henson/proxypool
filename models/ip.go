package models

import "gopkg.in/mgo.v2/bson"

// IP struct
type IP struct {
	ID   bson.ObjectId `bson:"_id" json:"-"`
	Data string        `bson:"data" json:"ip"`
	Type string        `bson:"type" json:"type"`
}

// NewIP .
func NewIP() *IP {
	return &IP{
		ID: bson.NewObjectId(),
	}
}
