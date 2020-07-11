package entity

import (
	"time"

	"github.com/goonode/mogo"
)

type RefreshSession struct {
	mogo.DocumentModel `bson:",inline" coll:"refresh_sessions"`
	Refresh            string `idx:"{refresh},unique" json:"refresh" binding:"required"`
	UserId             uint32 `json:"user_id" binding:"required"`
	ExpiresAt          time.Time `bson:"expires_at"`
}


func init() {
	mogo.ModelRegistry.Register(RefreshSession{})
}