package models

import (
	"belli/onki-game-ideas-mongo-backend/responses"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSpecData struct {
	ID             string
	Username       string
	Firstname      string
	Lastname       string
	HashedPassword string
	Role           responses.UserRole
	Status         responses.UserStatus
}

type UserLightData struct {
	ID        string
	Role      responses.UserRole
	Firstname string
	Lastname  string
}

type UserList struct {
	Total  int            `json:"total"`
	Result []UserSpecData `json:"result"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

/////////////////MONGO/////////////////////////////////////////

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
