package gorm

import (
	"github.com/Luke256/ducks/migration"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB, doMigration bool) (repo *GormRepository, init bool, err error) {
		if db == nil {
		return nil, false, gorm.ErrInvalidDB
	}
	repo = &GormRepository{
		db: db,
	}

	if doMigration {
		if init, err = migration.Migrate(db); err != nil {
			return nil, false, err
		}
	}

	return
}