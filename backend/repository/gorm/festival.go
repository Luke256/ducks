package gorm

import (
	"context"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GormRepository) RegisterFestival(name string, description string, date time.Time) (model.Festival, error) {
	festivalID, err := uuid.NewV7()
	if err != nil {
		return model.Festival{}, wrapGormError(err)
	}

	festival := model.Festival{
		ID:          festivalID,
		Name:        name,
		Description: description,
		Date:        date,
	}

	ctx := context.Background()

	if err := gorm.G[model.Festival](r.db).Create(ctx, &festival); err != nil {
		return model.Festival{}, wrapGormError(err)
	}

	return festival, nil
}

func (r *GormRepository) GetFestivalByID(festivalID uuid.UUID) (model.Festival, error) {
	ctx := context.Background()

	festival, err := gorm.G[model.Festival](r.db).
		Where(&model.Festival{ID: festivalID}, "ID").
		First(ctx)
	if err != nil {
		return model.Festival{}, wrapGormError(err)
	}

	return festival, nil
}

func (r *GormRepository) GetAllFestivals() ([]model.Festival, error) {
	ctx := context.Background()
	festivals, err := gorm.G[model.Festival](r.db).
		Order("date IS NULL ASC").Order("date ASC").Order("id ASC").Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}
	return festivals, nil
}

func (r *GormRepository) UpdateFestival(festivalID uuid.UUID, name string, description string, date time.Time) error {
	return r.mutateUnlockedFestivals(r.db.Model(&model.Festival{}).Select("id").Where("id = ?", festivalID), func(tx *gorm.DB) error {
		ctx := context.Background()

		rows, err := gorm.G[model.Festival](tx).
			Where(&model.Festival{ID: festivalID}, "ID").
			Select("Name", "Description", "Date").
			Updates(ctx, model.Festival{
				ID:          festivalID,
				Name:        name,
				Description: description,
				Date:        date,
			})
		if err != nil {
			return wrapGormError(err)
		}

		if rows == 0 {
			return wrapGormError(tx.Select("id").First(&model.Festival{}, "id = ?", festivalID).Error)
		}

		return nil
	})
}

func (r *GormRepository) DeleteFestival(festivalID uuid.UUID) error {
	return r.mutateUnlockedFestivals(r.db.Model(&model.Festival{}).Select("id").Where("id = ?", festivalID), func(tx *gorm.DB) error {
		ctx := context.Background()

		rows, err := gorm.G[model.Festival](tx).
			Where(&model.Festival{ID: festivalID}, "ID").
			Delete(ctx)
		if err != nil {
			return wrapGormError(err)
		}

		if rows == 0 {
			return repository.ErrNotFound
		}

		return nil
	})
}

func (r *GormRepository) SetFestivalLocked(festivalID uuid.UUID, locked bool) error {
	// ロック解除もできるよう、通常の変更禁止判定は適用しない。
	rows, err := gorm.G[model.Festival](r.db).
		Where("id = ?", festivalID).
		Update(context.Background(), "Locked", locked)
	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return wrapGormError(r.db.Select("id").First(&model.Festival{}, "id = ?", festivalID).Error)
	}
	return nil
}
