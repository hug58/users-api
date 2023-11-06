package token

import (
	"context"
)

// Repository handle the CRUD operations with Users.
type Repository interface {
	GetAll(ctx context.Context) ([]Token, error)
	Create(ctx context.Context, user_id uint, access string) error
	DeleteByUserId(ctx context.Context, user_id uint) error
}
