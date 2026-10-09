package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// v3 ポスター画像の複数枚対応
func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3",
		Migrate: func(db *gorm.DB) error {
			if err := db.AutoMigrate(&model.Poster{}, &model.PosterImage{}); err != nil {
				return err
			}
			if !db.Migrator().HasColumn(&model.Poster{}, "image_id") {
				return nil
			}

			// DDLが途中で失敗した後の再実行でも、移行済みの画像を重複登録しない。
			if err := db.Exec(`
				INSERT INTO poster_images (id, poster_id)
				SELECT p.image_id, p.id FROM posters AS p
				WHERE p.image_id <> '' AND NOT EXISTS (
					SELECT 1 FROM poster_images AS i
					WHERE i.id = p.image_id AND i.poster_id = p.id
				)
			`).Error; err != nil {
				return err
			}

			return db.Migrator().DropColumn(&model.Poster{}, "image_id")
		},
	}
}
