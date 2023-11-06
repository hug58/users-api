package token

import "time"

type Token struct {
	ID             int
	UserID         uint
	TokenValue     string
	ExpirationDate time.Time
	CreatedAt      time.Time
}
