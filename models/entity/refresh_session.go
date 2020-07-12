package entity

import (
	"time"
)

type RefreshSession struct {
	Refresh            string `bson:"refresh" json:"refresh"`
	UserEmail          string `bson:"user_email" json:"user_email"`
	CreatedAt          time.Time `bson:"created_at" json:"created_at"` 
}

func (obj *RefreshSession) New() *RefreshSession {
	return &RefreshSession{
		Refresh:     obj.Refresh,
		UserEmail:   obj.UserEmail,
		CreatedAt:   obj.CreatedAt,
	}
}
