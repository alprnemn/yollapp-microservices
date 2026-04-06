package repository

import (
	"context"
	"database/sql"
	"fmt"
	m "github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
)

type AuthRepository interface {
	CreateUserInvitation(ctx context.Context, inv *m.Invitation) error
	DeleteInvitation(ctx context.Context, userID string) error
	GetInvitationByToken(ctx context.Context, token string) (*m.Invitation, error)
}

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreateUserInvitation(ctx context.Context, inv *m.Invitation) error {

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

func (r *Repository) DeleteInvitation(ctx context.Context, userID string) error {

	query := `
	DELETE FROM user_invitations
	WHERE user_id = ?
	`

	result, err := r.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error executing query delete invitation: %s", err.Error())
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

func (r *Repository) GetInvitationByToken(ctx context.Context, token string) (*m.Invitation, error) {

	query := `SELECT * FROM user_invitations WHERE token = ?`

	inv := &m.Invitation{}

	err := r.DB.QueryRowContext(ctx, query, token).Scan(
		&inv.ID,
		&inv.UserID,
		&inv.Token,
		&inv.ExpiresAt,
		&inv.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("run query error getting invitation by token: %s", err.Error())
	}

	return inv, nil
}
