package services

import (
	"errors"
	"auth-task/config/db"
	"auth-task/models/entity"
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
