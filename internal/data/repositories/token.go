package repositories

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	pkgToken "github.com/hug58/users-api/pkg/token"
	"github.com/lib/pq"
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

func (ur *TokenRepository) DeleteByUserId(ctx context.Context, user_id uint) error {
	var query = `DELETE FROM tokens WHERE user_id=$1;`

	stmt, err := ur.Data.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, user_id); err != nil {
		pgErr, ok := err.(*pq.Error)

		if ok {
			if strings.Contains(pgErr.Message, " does not exist") {
				return nil
			}
			return err
		}

	}
	return nil
}
