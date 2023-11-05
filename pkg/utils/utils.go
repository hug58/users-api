package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	pkgUser "github.com/hug58/users-api/pkg/user"
)

type Message struct {
	Msg    string `json:"message"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Login struct {
	User        *pkgUser.User
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

func (p *Login) GenerarToken() error {
	//Create token
	token := jwt.New(jwt.SigningMethodHS256)

	if p.User == nil {
		return fmt.Errorf("user is required")
	}
	// Define  claims (data) the token
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = p.User.ID
	claims["phone"] = p.User.Phone
	claims["email"] = p.User.Email
	secret := os.Getenv("SECRET")

	tokenFirmed, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	p.AccessToken = string(tokenFirmed)
	p.TokenType = "bearer"

	return nil
}
