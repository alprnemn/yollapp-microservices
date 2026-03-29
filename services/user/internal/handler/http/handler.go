package http

import (
	"github.com/alprnemn/yollapp-microservices/services/user/internal/model"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/service"
	"github.com/alprnemn/yollapp-microservices/shared/utils"
	"log"
	"net/http"
	"strconv"
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
	router.HandleFunc("POST /user/create", h.HealthCheckHandler)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handler layer")

	id := req.PathValue("id")
	idInt, _ := strconv.Atoi(id)

	log.Println("id: ", idInt)

	user, err := h.Service.GetUser(req.Context(), idInt)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "ok",
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}

}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, req *http.Request) {

	var CreateUserDTO *model.CreateUserDTO

	if err := utils.ParseJSON(w, req, &CreateUserDTO); err != nil {
		utils.BadRequestResponse(w, req, err)
		return
	}


	h.Service.Create(>)

}
