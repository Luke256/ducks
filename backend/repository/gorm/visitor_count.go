package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func (r *GormRepository) ListVisitorCounts(ctx context.Context, festivalID uuid.UUID) ([]model.VisitorCount, error) {
	counts, err := gorm.G[model.VisitorCount](r.db).
		Where(&model.VisitorCount{FestivalID: festivalID}, "FestivalID").
		Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}
	return counts, nil
}

func (r *GormRepository) AddVisitorCount(ctx context.Context, festivalID uuid.UUID, timestamp time.Time, amount int) error {
	fes, err := gorm.G[model.Festival](r.db).
		Where(&model.Festival{ID: festivalID}, "ID").
		First(ctx)
	if err != nil {
		return wrapGormError(err)
	}

	// 10分ごとのバケットに丸める
	bucketStart := timestamp.UTC().Truncate(10 * time.Minute)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		count, err := gorm.G[model.VisitorCount](tx).
			Where(&model.VisitorCount{FestivalID: festivalID, BucketStart: bucketStart}, "FestivalID", "BucketStart").
			First(ctx)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return wrapGormError(err)
			}

			if amount < 0 {
				amount = 0
			}

			// レコードが存在しない場合は新規作成
			newCount := model.VisitorCount{
				FestivalID:  festivalID,
				BucketStart: bucketStart,
				Count:       uint(amount),
				Festival:    fes,
			}

			if err := gorm.G[model.VisitorCount](tx).Create(ctx, &newCount); err != nil {
				return wrapGormError(err)
			}
		} else {
			if amount < 0 && count.Count < uint(-amount) {
				count.Count = 0
			} else {
				count.Count += uint(amount)
			}

			// use Update to update the count
			_, err = gorm.G[model.VisitorCount](tx).
				Where(&model.VisitorCount{FestivalID: festivalID, BucketStart: bucketStart}, "FestivalID", "BucketStart").
				Update(ctx, "Count", count.Count)
			if err != nil {
				return wrapGormError(err)
			}
		}

		return nil
	})

	return wrapGormError(err)
}
