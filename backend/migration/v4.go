package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v4 販売商品の非表示機能
func v4() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "4",
		Migrate: func(db *gorm.DB) error {
			if !db.Migrator().HasColumn(&model.FestivalStock{}, "for_sale") {
				// 列の追加と既存商品の販売中への設定を一度に行う。
				if err := db.Exec("ALTER TABLE festival_stocks ADD COLUMN for_sale BOOLEAN NOT NULL DEFAULT TRUE").Error; err != nil {
					return err
				}
			}
			return db.AutoMigrate(&model.FestivalStock{})
		},
	}
}
