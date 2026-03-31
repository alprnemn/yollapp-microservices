package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	m "github.com/alprnemn/yollapp-microservices/services/user/internal/model"
	"github.com/alprnemn/yollapp-microservices/shared/errs"
	"github.com/lib/pq"
	"log"
	"time"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*m.User, error)
	Create(ctx context.Context, user *m.CreateUserDTO) error
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

func (r *Repository) GetByEmail(ctx context.Context, email string) (*m.User, error) {

	query := "SELECT id,first_name, last_name,username, email FROM users WHERE email = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &m.User{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		log.Println("ERR: ", err.Error())
		return nil, errs.ErrNotFound
	}

	return user, nil

}

func (r *Repository) Create(ctx context.Context, user *m.CreateUserDTO) error {

	query := `INSERT INTO users
			(first_name, last_name, username, email, phone, password)
			VALUES ($1,$2,$3,$4,$5,$6)
			`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := r.DB.ExecContext(ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Phone,
		user.Password,
	)
	if err != nil {
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) {
			switch pqErr.Constraint {
			case "users_username_key":
				return errs.ErrDuplicateUsername
			case "users_email_key":
				return errs.ErrDuplicateEmail
			case "users_phone_key":
				return errs.ErrDuplicatePhone
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("failed to create user: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}
