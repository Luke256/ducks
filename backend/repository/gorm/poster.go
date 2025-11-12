package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func (r *GormRepository) RegisterPoster(festivalID uuid.UUID, posterName, description, imageID string) (model.Poster, error) {
	posterID, err := uuid.NewV7()
	if err != nil {
		return model.Poster{}, err
	}

	var poster = model.Poster{
		ID:          posterID,
		FestivalID:  festivalID,
		PosterName:  posterName,
		Description: description,
		ImageID:     imageID,
		Status:      "uncollected",
	}

	ctx := context.Background()
	if err := gorm.G[model.Poster](r.db).Create(ctx, &poster); err != nil {
		return model.Poster{}, wrapGormError(err)
	}

	return poster, nil
}

func (r *GormRepository) GetPostersByFestivalID(festivalID uuid.UUID) ([]model.Poster, error) {
	ctx := context.Background()
	posters, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{FestivalID: festivalID}, "FestivalID").
		Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}
	return posters, nil
}

func (r *GormRepository) GetPosterByID(posterID uuid.UUID) (model.Poster, error) {
	ctx := context.Background()

	poster, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		First(ctx)
	if err != nil {
		return model.Poster{}, wrapGormError(err)
	}

	return poster, nil
}

func (r *GormRepository) GetPosterByFestivalIDAndPosterName(festivalID uuid.UUID, posterName string) (model.Poster, error) {
	ctx := context.Background()
	poster, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{FestivalID: festivalID, PosterName: posterName}, "FestivalID", "PosterName").
		First(ctx)
	if err != nil {
		return model.Poster{}, wrapGormError(err)
	}
	return poster, nil
}

func (r *GormRepository) UpdatePoster(posterID uuid.UUID, posterName, description string) error {
	ctx := context.Background()
	rows, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		Select("PosterName", "Description").
		Updates(ctx, model.Poster{PosterName: posterName, Description: description})
	if rows == 0 {
		return repository.ErrNotFound
	}
	return wrapGormError(err)
}

func (r *GormRepository) UpdatePosterStatus(posterID uuid.UUID, status string) error {
	ctx := context.Background()
	rows, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		Updates(ctx, model.Poster{Status: status})
	if rows == 0 {
		return repository.ErrNotFound
	}
	return wrapGormError(err)
}

func (r *GormRepository) DeletePoster(posterID uuid.UUID) error {
	ctx := context.Background()
	rowsAffected, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		Delete(ctx)
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return wrapGormError(err)
}
