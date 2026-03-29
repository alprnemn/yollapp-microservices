package service

import (
	"context"
	m "github.com/alprnemn/yollapp-microservices/services/user/internal/model"
	r "github.com/alprnemn/yollapp-microservices/services/user/internal/repository"
	"log"
)

type Service struct {
	repository r.UserRepository
}

func NewService(repo r.UserRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) GetUser(ctx context.Context, ID int) (*m.User, error) {
	log.Println("servicelayer")

	user, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *Service) Create(ctx context.Context, user m.CreateUserDTO) (int, error) {

	user, err := s.repository.Create(ctx, ID)
	if err != nil {
		return nil, err
	}

	return 0, nil

}
