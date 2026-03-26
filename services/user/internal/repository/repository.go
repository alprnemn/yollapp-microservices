package repository

import (
	"context"
	"database/sql"
	"github.com/alprnemn/yollapp-microservices/services/user/internal/model"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]model.User, error)
	RegisterUser(ctx context.Context, user model.User) (model.User, error)
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (r *Repository) RegisterUser(ctx context.Context) (*model.User, error) {
	return nil, nil
}
