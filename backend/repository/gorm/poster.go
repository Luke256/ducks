package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepository) RegisterPoster(festivalID uuid.UUID, posterName, description string, imageIDs []string) (model.Poster, error) {
	if len(imageIDs) > repository.MaxPosterImages {
		return model.Poster{}, repository.ErrInvalidPosterImages
	}
	posterID, err := uuid.NewV7()
	if err != nil {
		return model.Poster{}, err
	}

	var poster = model.Poster{
		ID:          posterID,
		FestivalID:  festivalID,
		PosterName:  posterName,
		Description: description,
		Status:      "uncollected",
		Images:      make([]model.PosterImage, len(imageIDs)),
	}
	for i, imageID := range imageIDs {
		poster.Images[i] = model.PosterImage{ID: imageID, PosterID: posterID}
	}

	ctx := context.Background()
	// 関連画像はINSERTし、既存の画像IDを別のポスターへ付け替えない。
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Images").Create(&poster).Error; err != nil {
			return err
		}
		if len(poster.Images) > 0 {
			return tx.Create(&poster.Images).Error
		}
		return nil
	}); err != nil {
		return model.Poster{}, wrapGormError(err)
	}

	return poster, nil
}

func (r *GormRepository) GetPostersByFestivalID(festivalID uuid.UUID) ([]model.Poster, error) {
	ctx := context.Background()
	posters, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{FestivalID: festivalID}, "FestivalID").
		Preload("Festival", nil).
		Preload("Images", nil).
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
		Preload("Festival", nil).
		Preload("Images", nil).
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
		Preload("Festival", nil).
		Preload("Images", nil).
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
	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return wrapGormError(r.db.Select("id").First(&model.Poster{}, "id = ?", posterID).Error)
	}
	return nil
}

func (r *GormRepository) UpdatePosterImages(posterID uuid.UUID, addIDs, deleteIDs []string) ([]model.PosterImage, error) {
	images := make([]model.PosterImage, 0)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 同じポスターの画像編集を直列化し、同時追加でも枚数上限を守る。
		var poster model.Poster
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&poster, "id = ?", posterID).Error; err != nil {
			return err
		}
		if err := tx.Where("poster_id = ?", posterID).Find(&images).Error; err != nil {
			return err
		}
		owned := make(map[string]bool, len(images))
		for _, img := range images {
			owned[img.ID] = true
		}
		for _, id := range deleteIDs {
			if !owned[id] {
				return repository.ErrInvalidPosterImages
			}
			delete(owned, id)
		}
		count := len(owned) + len(addIDs)
		if count < 1 || count > repository.MaxPosterImages {
			return repository.ErrInvalidPosterImages
		}
		if len(deleteIDs) > 0 {
			if err := tx.Where("poster_id = ? AND id IN ?", posterID, deleteIDs).Delete(&model.PosterImage{}).Error; err != nil {
				return err
			}
		}
		if len(addIDs) > 0 {
			added := make([]model.PosterImage, len(addIDs))
			for i, id := range addIDs {
				added[i] = model.PosterImage{ID: id, PosterID: posterID}
			}
			if err := tx.Create(&added).Error; err != nil {
				return err
			}
		}
		return tx.Where("poster_id = ?", posterID).Find(&images).Error
	})
	if err != nil {
		return nil, wrapGormError(err)
	}
	return images, nil
}

func (r *GormRepository) UpdatePosterStatus(posterID uuid.UUID, status string) error {
	ctx := context.Background()
	rows, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		Updates(ctx, model.Poster{Status: status})
	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return wrapGormError(r.db.Select("id").First(&model.Poster{}, "id = ?", posterID).Error)
	}
	return nil
}

func (r *GormRepository) DeletePoster(posterID uuid.UUID) error {
	ctx := context.Background()
	rows, err := gorm.G[model.Poster](r.db).
		Where(&model.Poster{ID: posterID}, "ID").
		Delete(ctx)

	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
