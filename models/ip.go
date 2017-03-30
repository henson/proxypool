package models

import "gopkg.in/mgo.v2/bson"

// IP struct
type IP struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Data string        `bson:"data" json:"ip"`
}

// NewIP .
func NewIP() *IP {
	return &IP{
		ID: bson.NewObjectId(),
	}
}
