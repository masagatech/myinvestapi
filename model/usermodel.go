package model

import "gopkg.in/mgo.v2/bson"

// User model
type User struct {
	Base
	Name     string        `json:"name" bson:"omitempty"`
	Password string        `json:"password" `
	Email    string        `json:"email"`
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
}
