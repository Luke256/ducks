package gorm

import (
	"gorm.io/gorm"
)

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		DB: db,
	}
}