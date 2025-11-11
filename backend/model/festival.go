package model

import (
	"github.com/google/uuid"
)

type Festival struct {
	ID          uuid.UUID `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
