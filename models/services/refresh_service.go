package services

import (
	"errors"
	"auth-task/config/db"
	"auth-task/models/entity"
	"time"
	"github.com/goonode/mogo"
	"labix.org/v2/mgo/bson"
)

type RefreshService struct{}

func (refreshservice RefreshService) Create(refreshSession *(entity.RefreshSession)) error {
	conn := db.GetConnection()
	defer conn.Session.Close()

	doc := mogo.NewDoc(entity.RefreshSession{}).(*(entity.RefreshSession))
	err := doc.FindOne(bson.M{"refresh": refreshSession.Refresh}, doc)

	if err == nil {
		return errors.New("Already Exist")
	}
	refreshSession.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	refreshModel := mogo.NewDoc(refreshSession).(*(entity.RefreshSession))
	err = mogo.Save(refreshModel)
	if vErr, ok := err.(*mogo.ValidationError); ok {
		return vErr
	}
	return err
}