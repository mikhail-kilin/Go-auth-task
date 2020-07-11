package entity

import (
	"github.com/dgrijalva/jwt-go"
	"auth-task/helpers"
	"time"
)

type User struct {
	Email              string `bson:"email" json:"email"`
	Password           string `bson:"password" json:"password"`
	Name               string `bson:"name" json:"name"`
}

func (u *User) New() *User {
	return &User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (user *User) GetJwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": string(user.Email),
		"name": string(user.Name),
		"created_at": time.Now(),
	})

	secretKey := helpers.EnvVar("SECRET")

	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}
