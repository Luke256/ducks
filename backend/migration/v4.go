package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v4 販売商品の非表示機能
func v4() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&model.FestivalStock{})
		},
	}
}
