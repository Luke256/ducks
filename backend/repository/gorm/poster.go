package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func (r *GormRepository) RegisterPoster(festivalID uuid.UUID, posterName string, description string, imageID string) (string, error) {
	var poster = model.Poster{
		FestivalID:  festivalID,
		PosterName:  posterName,
		Description: description,
		ImageID:     imageID,
		Status:      repository.PosterStatusUnCollected,
	}

	ctx := context.Background()
	err := gorm.G[model.Poster](r.db).Create(ctx, &poster)
	if err != nil {
		return "", err
	}

	return poster.ID.String(), nil
}

func (r *GormRepository) GetPostersByFestivalID(festivalID uuid.UUID) ([]model.Poster, error) {
	ctx := context.Background()
	posters, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{FestivalID: festivalID}).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	return posters, nil
}

func (r *GormRepository) GetPosterByID(posterID uuid.UUID) (model.Poster, error) {
	ctx := context.Background()

	poster, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}).
		First(ctx)
	if err != nil {
		return model.Poster{}, err
	}

	return poster, nil
}

func (r *GormRepository) GetPosterByFestivalIDAndPosterName(festivalID uuid.UUID, posterName string) (model.Poster, error) {
	ctx := context.Background()
	poster, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{FestivalID: festivalID, PosterName: posterName}).
		First(ctx)
	if err != nil {
		return model.Poster{}, err
	}
	return poster, nil
}

func (r *GormRepository) UpdatePosterStatus(posterID uuid.UUID, status string) error {
	ctx := context.Background()
	_, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}).
		Updates(ctx, model.Poster{Status: status})
	return wrapGormError(err)
}

func (r *GormRepository) DeletePoster(posterID uuid.UUID) error {
	ctx := context.Background()
	rowsAffected, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}).
		Delete(ctx)
	if rowsAffected == 0 {
		return repository.ErrPosterNotFound
	}

	return wrapGormError(err)
}
