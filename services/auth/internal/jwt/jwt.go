package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Authenticator struct {
	Secret   string
	Audience string
	Exp      time.Duration
}

func New(secret string, audience string, exp time.Duration) *Authenticator {
	return &Authenticator{
		Secret:   secret,
		Audience: audience,
		Exp:      exp,
	}
}

func (a *Authenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("error occurred while generating token")
	}
	return tokenString, nil
}
