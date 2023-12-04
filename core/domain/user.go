package domain

import "time"

type User struct {
	ID          string     `json:"id"`
	FullName    string     `json:"full_name"`
	PhoneNumber string     `json:"phone_number"`
	LoginCount  int        `json:"login_count"`
	Password    string     `json:"password"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type AuthData struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}
