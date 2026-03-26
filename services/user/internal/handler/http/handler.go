package http

import (
	"github.com/alprnemn/yollapp-microservices/pkg/utils"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/service"
	"net/http"
)

type Handler struct {
	Service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /user/{id}", h.GetUserHandler)
	router.HandleFunc("GET /user/health", h.HealthCheckHandler)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, req *http.Request) {

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{
		"user": "alprnemn",
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "ok",
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}

}
