package model


import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	Phone     string `db:"phone_number"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
