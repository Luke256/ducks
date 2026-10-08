package model

import (
	"time"

	"github.com/google/uuid"
)

type VisitorCount struct {
	FestivalID uuid.UUID `gorm:"type:char(36);not null;index;" json:"festival_id"`
	BucketStart time.Time `gorm:"" json:"bucket_start"`
	Count uint `json:"count"`

	Festival Festival `gorm:"foreignKey:FestivalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"festival"`
}