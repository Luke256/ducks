package model

import (
	"time"

	"github.com/google/uuid"
)

type Festival struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Locked      bool      `gorm:"not null;default:false" json:"locked"`
}
