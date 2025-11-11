package model

import (
	"github.com/google/uuid"
)

type Festival struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
