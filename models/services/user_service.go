package services

import (
	"errors"
	//"fmt"
	"auth-task/config/db"
	"auth-task/models/entity"

	"github.com/goonode/mogo"
	"labix.org/v2/mgo/bson"
)

type Userservice struct{}

func (userservice Userservice) Create(user *(entity.User)) error {
	conn := db.GetConnection()
	defer conn.Session.Close()
	
	doc := mogo.NewDoc(entity.User{}).(*(entity.User))
	err := doc.FindOne(bson.M{"email": user.Email}, doc)

	if err == nil {
		return errors.New("Already Exist")
	}
	userModel := mogo.NewDoc(user).(*(entity.User))
	err = mogo.Save(userModel)
	if vErr, ok := err.(*mogo.ValidationError); ok {
		return vErr
	}
	return err
}

/*
func (userservice Userservice) Find(user *(entity.User)) (*entity.User, error) {
	conn := db.GetConnection()
	defer conn.Session.Close()

	doc := mogo.NewDoc(entity.User{}).(*(entity.User))
	err := doc.FindOne(bson.M{"email": user.Email}, doc)

	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (userservice Userservice) FindByEmail(email string) (*entity.User, error) {
	conn := db.GetConnection()
	defer conn.Session.Close()

	user := new(entity.User)
	user.Email = email
	return userservice.Find(user)
}
*/
