package model

import (
	"github.com/google/uuid"
)

type StockItem struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Category    string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;"`
	ImageID     string    `gorm:"type:text;not null"`
}
