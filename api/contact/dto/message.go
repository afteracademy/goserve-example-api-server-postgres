package dto

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Msg       string    `json:"msg" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}
