package services

import (
	"context"
	"errors"
	//"fmt"
	"auth-task/config/db"
	"auth-task/models/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Userservice struct{}

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
/*
func (userservice Userservice) FindByEmail(email string) (*entity.User, error) {
	conn := db.GetConnection()
	defer conn.Session.Close()

	user := new(entity.User)
	user.Email = email
	return userservice.Find(user)
}
*/
