package model

type CreateUserResponseDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ActivateResponse struct {
	Message string `json:"message"`
}
