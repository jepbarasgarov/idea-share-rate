package models

import (
	"belli/onki-game-ideas-mongo-backend/responses"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserCreate struct {
	Username  string
	Firstname string
	Lastname  string
	Password  string
	Role      responses.UserRole
}

type UserUpdate struct {
	ID        primitive.ObjectID
	Username  string
	Firstname string
	Lastname  string
	Password  string
	Role      responses.UserRole
	Status    responses.UserStatus
}

type UserSpecDataBson struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	Username       string               `bson:"username" json:"username"`
	Firstname      string               `bson:"firstname" json:"firstname"`
	Lastname       string               `bson:"lastname" json:"lastname"`
	HashedPassword string               `bson:"password,omitempty" json:"password,omitempty"`
	Role           responses.UserRole   `bson:"role" json:"role"`
	Status         responses.UserStatus `bson:"status" json:"status"`
}

type UserLightDataBson struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Firstname string             `bson:"firstname" json:"firstname"`
	Lastname  string             `bson:"lastname" json:"lastname"`
	Role      responses.UserRole `bson:"role" json:"role"`
}
