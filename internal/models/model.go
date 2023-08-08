package models

import (
	"errors"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type User struct {
	ID        string     `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Name      string     `db:"name" json:"name"`
	Password  string     `db:"password" json:"password"`
	Role      string     `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type UserView struct {
	Username  string    `db:"username"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

type UserFilter struct {
	Name string `db:"name" json:"name"`
}

type UpdateName struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
