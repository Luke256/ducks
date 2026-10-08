package gorm

import (
	"time"

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
	// ドライバーの日時精度を保ったまま、自動生成時刻を UTC に揃える。
	now := db.NowFunc
	db = db.Session(&gorm.Session{
		NowFunc: func() time.Time { return now().UTC() },
	})
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
