package repositories

import (
	"context"
	"database/sql"
	"os"
	"time"

	pkgToken "github.com/hug58/users-api/pkg/token"
)

type TokenRepository struct {
	Data *sql.DB
}

func (ur *TokenRepository) GetAll(ctx context.Context) ([]pkgToken.Token, error) {
	return []pkgToken.Token{}, nil
}

func (ur *TokenRepository) Create(ctx context.Context, user_id uint, access string) error {
	query := `INSERT INTO tokens (user_id, token_value, expiration_date, created_at) 
			  VALUES ($1, $2, $3, $4) RETURNING id`

	duration, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRATION_MINUTES"))
	if err != nil {
		return err
	}

	token := pkgToken.Token{
		UserID:         user_id,
		TokenValue:     access,
		ExpirationDate: time.Now().Add(duration),
		CreatedAt:      time.Now(),
	}

	if err := ur.Data.QueryRowContext(ctx, query, token.UserID, token.TokenValue, token.ExpirationDate, token.CreatedAt).Scan(&token.ID); err != nil {
		return err
	}

	return nil
}
