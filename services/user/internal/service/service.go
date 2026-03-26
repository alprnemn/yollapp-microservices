package service

import (
	"context"
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

func (s *Service) GetUserByID(ctx context.Context, ID int) error {
	return nil
}
