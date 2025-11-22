package model

import (
	"time"

	"github.com/google/uuid"
)

type SaleRecord struct {
	ID              uuid.UUID `gorm:"type:string;primary_key"`
	FestivalStockID uuid.UUID `gorm:"type:string;not null;index"`
	Quantity        int       `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`

	FestivalStock FestivalStock `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
