package services

import (
	"context"
	"errors"
	"auth-task/config/db"
	"auth-task/models/entity"
	"auth-task/helpers"

	"github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"time"
)

type Userservice struct{}

type Tokens struct{
	Access_token   string
	Refresh_token  string
}

func (userservice Userservice) Create(user *(entity.User)) error {
	users_collection := db.GetConnection().DB.Collection("users")

	res, err := userservice.FindUser(user)

	if (res != nil ) {
		return errors.New("Already Exist")
	}

	ctx := context.Background()

	result, err := users_collection.InsertOne(ctx, user)

	if (err != nil || result == nil) {
		return errors.New("Something is wrong")
	}

	return nil
}

func (userservice Userservice) FindUser(info *entity.User) (*entity.User, error) {
	var user entity.User
	err := db.GetConnection().DB.Collection("users").FindOne(
		context.Background(),
		bson.M{"email": info.Email},
	).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}
	return &user, nil
}

func (userservice Userservice) GetTokens(user *entity.User) (Tokens, error) {
	current_time := time.Now()

	refresh_service := RefreshService{}
	session_id, refresh, errs := refresh_service.Generate(user, current_time)

	if errs != nil {
		return Tokens{}, errs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": string(user.Email),
		"name": string(user.Name),
		"created_at": current_time,
		"session_id": session_id,
	})

	secretKey := helpers.EnvVar("SECRET")

	tokenString, err := token.SignedString([]byte(secretKey))
	return Tokens{tokenString, refresh}, err
}

func (userservice Userservice) ReGenerateToken(access_token string) (Tokens, error) {
	secretKey := helpers.EnvVar("SECRET")

	token, err := jwt.Parse(access_token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return Tokens{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		session_id := claims["session_id"].(string)
		refresh_service := RefreshService{}

		refresh_service.DeleteSession(session_id)

		email := claims["email"].(string)
		user, err := userservice.FindUser(&entity.User {Email: email})
		if err != nil {
			return Tokens{}, err
		}

		return userservice.GetTokens(user)
	} else {
		return Tokens{}, errors.New("Something is wrong")
	}
}
