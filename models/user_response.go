package models

type ResponseUserStruct struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type UserProfileStruct struct {
	ID             int64   `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	Status         *string  `json:"status"`
	ProfilePicture *string `json:"profile_picture"`
	Role           string  `json:"role"`
}
