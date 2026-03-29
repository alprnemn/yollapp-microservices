package model

type User struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Age       uint8  `json:"age,omitempty"`
	Password  string `json:"-"`
}

type CreateUserDTO struct {
	User
}
