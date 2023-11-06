package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	pkgUser "github.com/hug58/users-api/pkg/user"
	"github.com/lib/pq"
)

type UserRepository struct {
	Data *sql.DB
}

func (ur *UserRepository) GetAll(ctx context.Context) ([]pkgUser.User, error) {
	var (
		users []pkgUser.User
	)

	query := `SELECT id, name, addres, email, phone,session_active, created_at, updated_at FROM users;`
	rows, err := ur.Data.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user pkgUser.User
		rows.Scan(&user.ID, &user.Name, &user.Address,
			&user.Email, &user.Phone, &user.SessionActive, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetOne(ctx context.Context, id uint) (pkgUser.User, error) {
	var (
		user  pkgUser.User
		query = `SELECT id, name, addres, email, phone, session_active, created_at, updated_at FROM users WHERE id = $1;`
	)

	row := ur.Data.QueryRowContext(ctx, query, id)
	if err := row.Scan(&user.ID, &user.Name, &user.Address, &user.Email, &user.Phone, &user.SessionActive,
		&user.CreatedAt, &user.UpdatedAt); err != nil {
		return pkgUser.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *pkgUser.User) error {
	var (
		query = `INSERT INTO users (name, addres, email, password, phone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	)

	if err := user.HashPassword(); err != nil {
		return err
	}

	row := ur.Data.QueryRowContext(
		ctx, query, user.Name, user.Address, user.Email,
		user.PasswordHash, user.Phone, time.Now(), time.Now(),
	)

	if err := row.Scan(&user.ID); err != nil {
		pgErr, ok := err.(*pq.Error)

		if ok {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_email_key") {
				return fmt.Errorf("email exists")
			} else if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_phone_key") {
				return fmt.Errorf("phone exists")
			}

			if pgErr.Code == "23514" && strings.Contains(pgErr.Message, "users_phone_check") {
				return fmt.Errorf("phone format not supported, please add a phone number as: code + number, example: +58123456999 ")
			}

		}
		return err
	}

	return nil
}

func (ur *UserRepository) Update(ctx context.Context, id uint, user pkgUser.User) (*pkgUser.User, error) {
	var (
		query = `UPDATE users SET name=$1, addres=$2, email=$3, phone=$4, updated_at=$5 WHERE id=$6;`
	)

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	stmt, err := ur.Data.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx, user.Name, user.Address, user.Email, user.Phone, time.Now(), id); err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	var query = `DELETE FROM users WHERE id=$1;`

	stmt, err := ur.Data.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Login(ctx context.Context, login *pkgUser.User) (*pkgUser.User, error) {

	var (
		user  pkgUser.User
		query = `SELECT id, name, addres, email, phone, password,created_at, updated_at FROM users WHERE phone = $1`
	)

	row := ur.Data.QueryRowContext(ctx, query, login.Phone)
	if err := row.Scan(&user.ID, &user.Name, &user.Address, &user.Email, &user.Phone, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	ok, err := user.PasswordMatch(login.Password)
	if err != nil {
		log.Println("testing")
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("password mismatch")
	}

	return &user, nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, id uint, password string) (string, error) {
	var (
		query = `UPDATE users SET password=$1, updated_at=$2 WHERE id=$3;`
		user  = pkgUser.User{
			Password: password,
		}
	)

	if err := user.HashPassword(); err != nil {
		return "", err
	}

	stmt, err := ur.Data.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}

	defer stmt.Close()
	if _, err := stmt.ExecContext(ctx, user.PasswordHash, time.Now(), id); err != nil {
		return "", err
	}

	return "Changed password Succesfully", nil

}
