package model

import "time"

type UserInvitation struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}
