package model

import (
	"github.com/google/uuid"
)

type Poster struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	FestivalID  uuid.UUID `gorm:"type:char(36);not null;index:idx_poster,priority:1;" json:"festival_id"`
	PosterName  string    `gorm:"type:char(64);not null;index:idx_poster,priority:2;" json:"poster_name"`
	Description string    `gorm:"type:text;not null;" json:"description"`
	ImageID     string    `gorm:"type:text;not null;" json:"image_id"`
	Status      string    `gorm:"type:text;not null;" json:"status"`

	Festival Festival `gorm:"foreignKey:FestivalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"festival"`
}
