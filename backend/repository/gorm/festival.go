package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GormRepository) RegisterFestival(name string, description string) (model.Festival, error) {
	festivalID, err := uuid.NewV7()
	if err != nil {
		return model.Festival{}, wrapGormError(err)
	}

	festival := model.Festival{
		ID:          festivalID,
		Name:        name,
		Description: description,
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
		Where(&model.Festival{ID: festivalID}).
		First(ctx)
	if err != nil {
		return model.Festival{}, wrapGormError(err)
	}

	return festival, nil
}

func (r *GormRepository) GetAllFestivals() ([]model.Festival, error) {
	ctx := context.Background()
	festivals, err := gorm.G[model.Festival](r.db).Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}
	return festivals, nil
}

func (r *GormRepository) UpdateFestival(festivalID uuid.UUID, name string, description string) error {
	ctx := context.Background()

	rows, err := gorm.G[model.Festival](r.db).
		Where(&model.Festival{ID: festivalID}).
		Select("Name", "Description").
		Updates(ctx, model.Festival{
			ID:          festivalID,
			Name:        name,
			Description: description,
		})
	if err != nil {
		return wrapGormError(err)
	}

	if rows == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormRepository) DeleteFestival(festivalID uuid.UUID) error {
	ctx := context.Background()

	rows, err := gorm.G[model.Festival](r.db).
		Where(&model.Festival{ID: festivalID}).
		Delete(ctx)
	if err != nil {
		return wrapGormError(err)
	}

	if rows == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}