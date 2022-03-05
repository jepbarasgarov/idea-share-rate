package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WorkerCreate struct {
	Firstname string
	Lastname  string
	Position  string
}

type WorkerUpdate struct {
	ID        primitive.ObjectID
	Firstname string
	Lastname  string
	Position  string
}

type WorkerLightData struct {
	ID        string
	Firstname string
	Lastname  string
	Position  string
}

//////////////////////////////////////////////////////////////////////BSON////////////////////////////////////////////////////////////////////////////////

type WorkerBsonModelInIdea struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Firstname string             `bson:"firstname,omitempty" json:"firstname"`
	LastName  string             `bson:"lastname,omitempty" json:"lastname"`
	Position  string             `bson:"position,omitempty" json:"position"`
}
