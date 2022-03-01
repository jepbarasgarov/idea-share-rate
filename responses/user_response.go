package responses

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

type UserList struct {
	Total  int            `json:"total"`
	Result []UserSpecData `json:"result"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ActionInfo struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	Language  Lang       `json:"language"`
}

type UserLogin struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Firstname    string     `json:"firstname"`
	Lastname     string     `json:"lastname"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}

type UserSpecData struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
}

type UserLightData struct {
	ID        string   `json:"id"`
	Role      UserRole `json:"role"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
}
