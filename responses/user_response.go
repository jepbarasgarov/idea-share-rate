package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserRole string

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
)

type UserStatus string

const (
	Active  UserStatus = "ACTIVE"
	Blocked UserStatus = "BLOCKED"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ActionInfo struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Firstname string             `json:"firstname"`
	Lastname  string             `json:"lastname"`
	Role      UserRole           `json:"role"`
	Status    UserStatus         `json:"status"`
	Language  Lang               `json:"language"`
}

type UserLogin struct {
	ID           primitive.ObjectID `json:"id"`
	Username     string             `json:"username"`
	Firstname    string             `json:"firstname"`
	Lastname     string             `json:"lastname"`
	Role         UserRole           `json:"role"`
	Status       UserStatus         `json:"status"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}
