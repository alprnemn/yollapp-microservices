package service

import (
	"context"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/gateway/http/user"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"github.com/alprnemn/yollapp-microservices/shared/errs"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userGateway *user.Gateway
}

func New(usergateway *user.Gateway) *Service {
	return &Service{
		userGateway: usergateway,
	}
}

func (s *Service) Login() error {
	return nil
}

func (s *Service) Register(ctx context.Context, payload *model.RegisterUserDTO) (*model.RegisterUserResponseDTO, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.ErrGeneratePassword
	}

	req := &model.RegisterUserDTO{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Password:  string(hashedPassword),
	}

	rsp, err := s.userGateway.RegisterUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *Service) ActivateUser() error {
	return nil
}

func (s *Service) RefreshToken() error {
	return nil
}

func (s *Service) ResetPassword() error {
	return nil
}
func (s *Service) Logout() error {
	return nil
}
