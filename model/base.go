package model

import "time"

// Base contains common fields for all tables
type Base struct {
	CreatedAt time.Time `json:"created_at" bson:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"omitempty"`
	DeletedAt time.Time `json:"deleted_at" bson:"omitempty"`
}

type Sysparams struct {
	Schema  string `json:"schema"`
	Operate string `json:"operate"`
	Flag    string `json:"flag"`
	Payload string `json:"payload"`
}
