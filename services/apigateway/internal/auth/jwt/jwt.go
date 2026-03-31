package jwt

import (
	"errors"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenValidator struct {
	secret string
	exp    time.Duration
	issuer string
}

func NewTokenValidator(exp time.Duration, secret, issuer string) *TokenValidator {
	return &TokenValidator{
		secret: secret,
		exp:    exp,
		issuer: issuer,
	}
}

//func (a *Authenticator) GenerateToken(claims jwt.Claims) (string, error) {
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	tokenString, err := token.SignedString([]byte(a.secret))
//	if err != nil {
//		return "", errs.New("error occurred while generating token")
//	}
//	return tokenString, nil
//}

func (t *TokenValidator) Validate(token string) (*jwt.Token, error) {

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(t.secret), nil
	}

	tk, err := jwt.Parse(token, keyFunc, jwt.WithExpirationRequired(),
		jwt.WithAudience(t.issuer),
		jwt.WithIssuer(t.issuer),
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
	if !tk.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return tk, nil
}
