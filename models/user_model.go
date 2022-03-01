package models

import (
	"belli/onki-game-ideas-mongo-backend/responses"
)

type UserCreate struct {
	Username  string
	Firstname string
	Lastname  string
	Password  string
	Role      responses.UserRole
}

type UserUpdate struct {
	ID        string
	Username  string
	Firstname string
	Lastname  string
	Password  string
	Role      responses.UserRole
	Status    responses.UserStatus
}

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
