package http

import "net/http"

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /auth/login", h.LoginHandler)
	router.HandleFunc("POST /auth/register", h.RegisterHandler)
	router.HandleFunc("GET /auth/activate", h.ActivateUserHandler)
	router.HandleFunc("POST /auth/refresh", h.RefreshTokenHandler)
	router.HandleFunc("POST /auth/resetpassword", h.ResetPasswordHandler)
	router.HandleFunc("POST /auth/logout", h.LogoutHandler)
}
