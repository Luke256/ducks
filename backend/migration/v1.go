package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v1 販売管理システムの追加
func v1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&model.StockItem{},
				&model.FestivalStock{},
				&model.SaleRecord{},
			)
		},
	}
}