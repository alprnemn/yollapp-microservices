package service

import (
	"context"
	m "github.com/alprnemn/yollapp-microservices/services/user/internal/model"
	r "github.com/alprnemn/yollapp-microservices/services/user/internal/repository"
)

type Service struct {
	repository r.UserRepository
}

func NewService(repo r.UserRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) Create(ctx context.Context, user *m.CreateUserDTO) (*m.CreateUserResponseDTO, error) {

	err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	usr, err := s.repository.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	resp := &m.CreateUserResponseDTO{
		ID:       usr.ID,
		Username: usr.Username,
		Email:    usr.Email,
	}

	return resp, nil
}

func (s *Service) Activate(ctx context.Context, user *m.ActivateUserDTO) error {

	err := s.repository.ActivateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
