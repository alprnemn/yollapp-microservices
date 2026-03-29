package http

import (
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/service"
	"github.com/alprnemn/yollapp-microservices/shared/utils"
	"log"
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
	var registerUserDTO *model.RegisterUserDTO

	if err := utils.ParseJSON(w, req, &registerUserDTO); err != nil {
		utils.BadRequestResponse(w, req, err)
		return
	}

	log.Println("username: ", registerUserDTO.Username)
	log.Println("password: ", registerUserDTO.Password)

}
func (h *Handler) ActivateUserHandler(w http.ResponseWriter, req *http.Request)  {}
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, req *http.Request)  {}
func (h *Handler) ResetPasswordHandler(w http.ResponseWriter, req *http.Request) {}
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request)        {}
