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

type UserService struct{}

type Tokens struct{
	AccessToken   string
	RefreshToken  string
}

func (userService UserService) Create(user *(entity.User)) error {
	usersCollection := db.GetConnection().DB.Collection("users")

	res, err := userService.FindUser(user)

	if (res != nil ) {
		return errors.New("Already Exist")
	}

	ctx := context.Background()

	result, err := usersCollection.InsertOne(ctx, user)

	if (err != nil || result == nil) {
		return errors.New("Something is wrong")
	}

	return nil
}

func (userService UserService) FindUser(info *entity.User) (*entity.User, error) {
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

func (userService UserService) GetTokens(user *entity.User) (Tokens, error) {
	currentTime := time.Now()

	refreshService := RefreshService{}
	sessionId, refresh, errs := refreshService.Generate(user, currentTime)

	if errs != nil {
		return Tokens{}, errs
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": string(user.Email),
		"name": string(user.Name),
		"created_at": currentTime,
		"session_id": sessionId,
	})

	secretKey := helpers.EnvVar("SECRET")

	tokenString, err := token.SignedString([]byte(secretKey))
	return Tokens{tokenString, refresh}, err
}

func (userService UserService) ReGenerateToken(accessToken string) (Tokens, error) {
	secretKey := helpers.EnvVar("SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return Tokens{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sessionId := claims["session_id"].(string)
		refreshService := RefreshService{}

		refreshService.DeleteSession(sessionId)

		email := claims["email"].(string)
		user, err := userService.FindUser(&entity.User {Email: email})
		if err != nil {
			return Tokens{}, err
		}

		return userService.GetTokens(user)
	} else {
		return Tokens{}, errors.New("Something is wrong")
	}
}
