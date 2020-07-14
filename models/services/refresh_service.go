package services

import (
	"crypto/sha512"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"auth-task/config/db"
	"auth-task/models/entity"
	"auth-task/helpers"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/dgrijalva/jwt-go"
)

type RefreshService struct{}

func (refreshservice RefreshService) Create(refreshSession *(entity.RefreshSession)) (string, error) {
	sessionCollection := db.GetConnection().DB.Collection("refresh_sessions")

	ctx := context.Background()
	refreshSession.CreatedAt = time.Now()

	result, err := sessionCollection.InsertOne(ctx, refreshSession)

	if (err != nil) {
		return "Error", errors.New("Something is wrong")
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}


func (refreshservice RefreshService) Generate(user *(entity.User), tokenTime time.Time) (string, string, error) {
	token := refreshservice.GetRefreshToken(user, tokenTime.String())

	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)

	if (err != nil) {
		return "error", "error", errors.New("Something is wrong")
	}

	session := entity.RefreshSession{string(hash), user.Email, time.Now()}
	sesssionId, errc := refreshservice.Create(&session)
	if (errc != nil) {
		return "error", "error", errors.New("Something is wrong")
	}


	return sesssionId, token, nil
}

func (refreshservice RefreshService) GetRefreshToken(user *(entity.User), tokenTime string) (string) {
	secret := helpers.EnvVar("SECRET")

	h512 := sha512.New()
	token := base64.StdEncoding.EncodeToString(h512.Sum([]byte(user.Email + user.Name + tokenTime + secret)))

	return token
}

func (refreshservice RefreshService) FindSession(id string) (*entity.RefreshSession, error) {
	var session entity.RefreshSession

	result, errh := primitive.ObjectIDFromHex(id)

	if (errh != nil) {
		return nil, errors.New("Something is wrong")
	}

	err := db.GetConnection().DB.Collection("refresh_sessions").FindOne(
		context.Background(),
		bson.M{"_id": result},
	).Decode(&session)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}
	return &session, nil
}

func (refreshservice RefreshService) DeleteSession(id string) (error) {
	result, errh := primitive.ObjectIDFromHex(id)

	if (errh != nil) {
		return errors.New("Something is wrong")
	}

	db.GetConnection().DB.Collection("refresh_sessions").DeleteOne(
		context.Background(),
		bson.M{"_id": result},
	)
	return nil
}

func (refreshservice RefreshService) DeleteSessionByAccessToken(accessToken string) (error) {
	defer db.CloseConection()

	secretKey := helpers.EnvVar("SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sessionId := claims["session_id"].(string)

		refreshservice.DeleteSession(sessionId)

		return nil
	} else {
		return errors.New("Something is wrong")
	}
}

func (refreshservice RefreshService) DeleteManySessionsByUserEmail(email string) error{
	sessionCollection := db.GetConnection().DB.Collection("refresh_sessions")

	ctx := context.Background()

	result, err := sessionCollection.DeleteMany(ctx, bson.M{"user_email": email})

	if (err != nil || result == nil) {
		return errors.New("Something is wrong")
	}

	return nil
}

func (refreshservice RefreshService) DeleteAllSessionsOfUser(accessToken string) (error) {
	defer db.CloseConection()

	secretKey := helpers.EnvVar("SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)

		refreshservice.DeleteManySessionsByUserEmail(email)

		return nil
	} else {
		return errors.New("Something is wrong")
	}
}
