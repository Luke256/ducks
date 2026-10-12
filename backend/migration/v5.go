package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v5 イベントのロック状態の追加
func v5() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&model.Festival{})
		},
	}
}
