package model

type CreateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type ActivateUserDTO struct {
	ID string `json:"userid"`
}
