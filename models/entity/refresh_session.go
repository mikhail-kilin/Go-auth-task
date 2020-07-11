package entity

import (
	"time"

	"github.com/goonode/mogo"
)

type RefreshSession struct {
	Refresh            string `bson:"refresh" json:"refresh"`
	UserEmail          string `bson:"user_email" json:"user_email"`
	ExpiresAt          time.Time `bson:"expires_at" json:"expires_at"`
}

func (obj *RefreshSession) New() *RefreshSession {
	return &RefreshSession{
		Refresh:     obj.Refresh,
		UserEmail:      obj.UserEmail,
		ExpiresAt:   obj.ExpiresAt,
	}
}

