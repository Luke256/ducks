package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GormRepository) RegisterStockItem(name string, description string, category string, imageID string) (model.StockItem, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return model.StockItem{}, err
	}

	var item = model.StockItem{
		ID:          id,
		Name:        name,
		Description: description,
		Category:    category,
		ImageID:     imageID,
	}

	ctx := context.Background()

	if err := gorm.G[model.StockItem](r.db).Create(ctx, &item); err != nil {
		return model.StockItem{}, wrapGormError(err)
	}

	return item, nil
}

func (r *GormRepository) GetStockItemByID(id uuid.UUID) (model.StockItem, error) {
	ctx := context.Background()

	item, err := gorm.G[model.StockItem](r.db).
		Where(model.StockItem{ID: id}, "ID").
		First(ctx)
	if err != nil {
		return model.StockItem{}, wrapGormError(err)
	}

	return item, nil
}

func (r *GormRepository) QueryStockItems(category string) ([]model.StockItem, error) {
	ctx := context.Background()

	var query = model.StockItem{
		Category: category,
	}

	items, err := gorm.G[model.StockItem](r.db).
		Where(query).
		Find(ctx)

	if err != nil {
		return nil, wrapGormError(err)
	}

	return items, nil
}

func (r *GormRepository) UpdateStockItem(id uuid.UUID, name string, description string, category string, imageID string) (model.StockItem, error) {
	ctx := context.Background()

	var item = model.StockItem{
		ID:          id,
		Name:        name,
		Description: description,
		Category:    category,
		ImageID:     imageID,
	}

	err := r.mutateUnlockedFestivals(r.db.Model(&model.FestivalStock{}).Select("festival_id").Where("stock_item_id = ?", id), func(tx *gorm.DB) error {
		rows, err := gorm.G[model.StockItem](tx).
			Where(item, "ID").
			Select("Name", "Description", "Category", "ImageID").
			Updates(ctx, item)
		if err != nil {
			return err
		}
		if rows == 0 {
			if err := tx.Select("id").First(&model.StockItem{}, "id = ?", id).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return model.StockItem{}, err
	}

	return item, nil
}

func (r *GormRepository) DeleteStockItem(id uuid.UUID) error {
	return r.mutateUnlockedFestivals(r.db.Model(&model.FestivalStock{}).Select("festival_id").Where("stock_item_id = ?", id), func(tx *gorm.DB) error {
		ctx := context.Background()

		rows, err := gorm.G[model.StockItem](tx).
			Where(model.StockItem{ID: id}, "ID").
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
