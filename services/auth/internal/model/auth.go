package model

import "time"

type RegisterUserDTO struct {
	FirstName string `json:"firstname" validate:"required,min=3,max=20"`
	LastName  string `json:"lastname" validate:"required,min=3,max=20"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email,max=55"`
	Phone     string `json:"phone" validate:"required,min=9,max=20"`
	Password  string `json:"password" validate:"required,min=6,max=25"`
}

type ActivateResponse struct {
	Message string `json:"message"`
}

type ActivateUserDTO struct {
	ID string `json:"userid"`
}

type Invitation struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
type CreateUserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserResponseDTO struct {
	CreateUserResponseDTO
	Token string `json:"token"`
}

func MapCreateUserResponseToRegisterResponse(user CreateUserResponseDTO, token string) *RegisterUserResponseDTO {
	return &RegisterUserResponseDTO{
		user,
		token,
	}
}
