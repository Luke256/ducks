package gorm

import (
	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 開催後にロックしたイベントへの変更を拒否し、判定と更新の間のロック状態変更を防ぐ。
func (r *GormRepository) mutateUnlockedFestivals(festivalIDs *gorm.DB, mutate func(*gorm.DB) error) error {
	var ids []uuid.UUID
	if err := festivalIDs.Find(&ids).Error; err != nil {
		return wrapGormError(err)
	}
	return wrapGormError(r.db.WithContext(festivalIDs.Statement.Context).Transaction(func(tx *gorm.DB) error {
		var festivals []model.Festival
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", ids).Order("id").Find(&festivals).Error; err != nil {
			return err
		}
		for _, festival := range festivals {
			if festival.Locked {
				return repository.ErrEventLocked
			}
		}
		return mutate(tx)
	}))
}
