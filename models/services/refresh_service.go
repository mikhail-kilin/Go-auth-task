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
)

type RefreshService struct{}

func (refreshservice RefreshService) Create(refreshSession *(entity.RefreshSession)) error {
	session_collection := db.GetConnection().DB.Collection("refresh_sessions")

	ctx := context.Background()
	refreshSession.CreatedAt = time.Now()

	result, err := session_collection.InsertOne(ctx, refreshSession)

	if (err != nil || result == nil) {
		return errors.New("Something is wrong")
	}

	return nil
}

func (refreshservice RefreshService) Generate(user *(entity.User), jwt string) (string, error) {
	secret := helpers.EnvVar("SECRET")

	h512 := sha512.New()
	token := base64.StdEncoding.EncodeToString(h512.Sum([]byte(jwt + secret)))

	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)

	if (err != nil) {
		return "error", errors.New("Something is wrong")
	}

	session := entity.RefreshSession{string(hash), user.Email, time.Now()}
	errc := refreshservice.Create(&session)
	if (errc != nil) {
		return "error", errors.New("Something is wrong")
	}


	return token, nil
}
