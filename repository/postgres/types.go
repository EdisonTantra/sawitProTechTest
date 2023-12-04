package postgres

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `db:"id" sql:",type:uuid"`
	FullName    string     `db:"full_name"`
	PhoneNumber string     `db:"phone_number"`
	Password    string     `db:"password"`
	LoginCount  int        `db:"login_count"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type UserPatchArg struct {
	ID          uuid.UUID `db:"id" sql:",type:uuid"`
	FullName    string    `db:"full_name"`
	PhoneNumber string    `db:"phone_number"`
}

type UserLoginArg struct {
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}
