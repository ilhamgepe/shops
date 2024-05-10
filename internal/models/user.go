package models

import "time"

type Role string

var (
	UserRole  Role = "users"
	AdminRole Role = "admin"
)

type CreateUserDTO struct {
	Firstname string  `db:"first_name" json:"firstname" validate:"required,min=2,max=20"`
	Lastname  *string `db:"last_name" json:"lastname" validate:"omitempty,min=2,max=20"`
	Email     string  `db:"email" json:"email" validate:"required,email"`
	Password  string  `db:"password" json:"password" validate:"required,min=6,max=100"`
}

type LoginUserDTO struct {
	Email    string `db:"email" json:"email" validate:"required,email"`
	Password string `db:"password" json:"password" validate:"required,min=6,max=100"`
}

type User struct {
	ID        int64      `db:"id" json:"id"`
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  *string    `db:"last_name" json:"last_name,omitempty"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	Role      Role       `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
