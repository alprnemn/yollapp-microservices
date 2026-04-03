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
	Role      Role
	Address   *Address
}

type Role struct {
	ID   int
	Name string
}
