package http

import (
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/service"
	"github.com/alprnemn/yollapp-microservices/services/auth/pkg/validator"
	"github.com/alprnemn/yollapp-microservices/shared/utils"
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

func (h *Handler) LoginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello from login handler")

}

func (h *Handler) RegisterHandler(w http.ResponseWriter, req *http.Request) {

	var payload model.RegisterUserDTO

	if err := utils.ParseJSON(w, req, &payload); err != nil {
		utils.BadRequestResponse(w, req, err)
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		utils.BadRequestResponse(w, req, err)
		return
	}

	ctx := req.Context()

	res, err := h.Service.Register(ctx, &payload)
	if err != nil {
		utils.BadRequestResponse(w, req, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, res); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error occccccuurreedd")
		return
	}

}
func (h *Handler) ActivateUserHandler(w http.ResponseWriter, req *http.Request)  {}
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, req *http.Request)  {}
func (h *Handler) ResetPasswordHandler(w http.ResponseWriter, req *http.Request) {}
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request)        {}
