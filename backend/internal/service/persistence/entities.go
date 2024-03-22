package persistence

import (
	"github.com/google/uuid"
	"time"
)

type (
	Entry struct {
		ID        uuid.UUID `gorm:"id,primaryKey" json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}
)
