package model

import (
	"github.com/google/uuid"
)

type FestivalStock struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	FestivalID  uuid.UUID `gorm:"type:char(36);not null;index"`
	StockItemID uuid.UUID `gorm:"type:char(36);not null;index"`
	Price       int       `gorm:"not null"`

	Festival Festival  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Stock    StockItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
