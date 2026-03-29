package repository

import (
	"context"
	"database/sql"
	"errors"
	m "github.com/alprnemn/yollapp-microservices/services/user/internal/model"
	"log"
	"time"
)

type UserRepository interface {
	GetByID(ctx context.Context, ID int) (*m.User, error)
	Create(ctx context.Context, user m.CreateUserDTO) (int, error)
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

const QueryTimeoutDuration = time.Second * 3

func (r *Repository) GetByID(ctx context.Context, ID int) (*m.User, error) {
	log.Println("repository layer")
	query := "SELECT id,first_name, last_name,username, email,age FROM users WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &m.User{}
	err := r.DB.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Age,
	)
	if err != nil {
		return nil, errors.New("error getting user")
	}

	return user, nil

}

func (r *Repository) Create(ctx context.Context, user m.CreateUserDTO) (int, error) {
	query := "INSERT INTO users VALUES ($1,$2,$3,$4,$5,$6,$7)"

	result, err := r.DB.ExecContext(ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Age,
	)
	if err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if n == 0 || err != nil {
		return 0, err
	}
}
