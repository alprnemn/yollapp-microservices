package jwt

import (
	"context"
	"errors"
	"fmt"
	u "github.com/alprnemn/yollapp-microservices/pkg/utils"
	e "github.com/alprnemn/yollapp-microservices/services/apigateway/pkg/errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// login

type Authenticator struct {
	secret string
	exp    time.Duration
	issuer string
}

func NewJWTAuthenticator(exp time.Duration, secret, issuer string) *Authenticator {
	return &Authenticator{
		secret: secret,
		exp:    exp,
		issuer: issuer,
	}
}

func (a *Authenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", errors.New("error occurred while generating token")
	}
	return tokenString, nil
}

func (a *Authenticator) ValidateToken(token string) (*jwt.Token, error) {

	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	}

	t, err := jwt.Parse(token, keyFunc, jwt.WithExpirationRequired(),
		jwt.WithAudience(a.issuer),
		jwt.WithIssuer(a.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {

		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, fmt.Errorf("malformed token")

		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, fmt.Errorf("invalid signature")

		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("token expired")

		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, fmt.Errorf("token not active yet")

		default:
			return nil, fmt.Errorf("could not validate token: %w", err)
		}
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return t, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO : Check after user and auth completed

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			u.BadRequestResponse(w, r, e.ErrMissingAuthHeader)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			u.UnauthorizedError(w, r, e.ErrWrongAuthHeader)
			return
		}

		token, err := a.ValidateToken(parts[1])
		if err != nil {
			u.UnauthorizedError(w, r, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
