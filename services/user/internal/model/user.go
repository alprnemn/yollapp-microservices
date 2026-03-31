package model

import (
	"time"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Username  string
	Email     string
	Age       int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Password  string
}

type CreateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type CreateUserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
