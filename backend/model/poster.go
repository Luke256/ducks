package model

import (
	"github.com/google/uuid"
)

type Poster struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	FestivalID  uuid.UUID `gorm:"type:char(36);not null;index:idx_poster,priority:1;" json:"festival_id"`
	PosterName  string    `gorm:"type:char(64);not null;index:idx_poster,priority:2;" json:"poster_name"`
	Description string    `gorm:"type:text;not null;" json:"description"`
	Status      string    `gorm:"type:text;not null;" json:"status"`

	Festival Festival      `gorm:"foreignKey:FestivalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"festival"`
	Images   []PosterImage `gorm:"foreignKey:PosterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"images"`
}

type PosterImage struct {
	ID       string    `gorm:"type:varchar(255);primary_key;" json:"id"`
	PosterID uuid.UUID `gorm:"type:char(36);not null;index:idx_poster_image,priority:1;" json:"poster_id"`
}
