package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/gateway/http/user"
	j "github.com/alprnemn/yollapp-microservices/services/auth/internal/jwt"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/mailer"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/model"
	"github.com/alprnemn/yollapp-microservices/services/auth/internal/repository"
	m "github.com/alprnemn/yollapp-microservices/services/auth/model"
	"github.com/alprnemn/yollapp-microservices/shared/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository    repository.AuthRepository
	userGateway   *user.Gateway
	Authenticator *j.Authenticator
	Mailer        mailer.Mailer
}

func New(usergateway *user.Gateway, authenticator *j.Authenticator, mailer mailer.Mailer, repo repository.AuthRepository) *Service {
	return &Service{
		userGateway:   usergateway,
		Authenticator: authenticator,
		Mailer:        mailer,
		Repository:    repo,
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

	claims := jwt.MapClaims{
		"sub": rsp.ID,
		"exp": time.Now().Add(s.Authenticator.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": s.Authenticator.Audience,
		"aud": s.Authenticator.Audience,
	}

	token, err := s.Authenticator.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	rsp.Token = token

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	inv := &m.UserInvitation{
		UserID:    rsp.ID,
		Token:     hashToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.Repository.CreateUserInvitation(ctx, inv)
	if err != nil {
		return nil, err
	}

	html := fmt.Sprintf("<strong> Hello there click the link to activate your account: <a>http://127.0.01:8080/auth/activate?token=%s</a></strong>", plainToken)

	if err := s.Mailer.Send(ctx, rsp.Email, "Create user", html, rsp.ID); err != nil {
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
