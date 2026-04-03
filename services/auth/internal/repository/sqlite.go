package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/model"
	"log"
)

type AuthRepository interface {
	CreateUserInvitation(ctx context.Context, inv *model.UserInvitation) error
	DeleteInvitation(ctx context.Context, token string) error
}

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreateUserInvitation(ctx context.Context, inv *model.UserInvitation) error {
	log.Println("hello from repo layer: ", inv.Token)
	query := `
	INSERT INTO user_invitations 
	(user_id, token, expires_at)
	VALUES (?, ?, ?)
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		inv.UserID,
		inv.Token,
		inv.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("error executing query: %s", err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", rows)
	}

	return nil
}

func (r *Repository) DeleteInvitation(ctx context.Context, token string) error {

	query := `
	DELETE FROM user_invitations
	WHERE token = ?
	`

	result, err := r.DB.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("error executing query: %s", err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("invitation not found")
	}

	return nil
}
