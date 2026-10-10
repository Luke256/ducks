package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

func (r *GormRepository) RegisterFestivalStock(festivalID, itemID uuid.UUID, price int, description string, forSale bool) (model.FestivalStock, error) {
	fesStockID, err := uuid.NewV7()
	if err != nil {
		return model.FestivalStock{}, err
	}

	var stock = model.FestivalStock{
		ID:          fesStockID,
		FestivalID:  festivalID,
		StockItemID: itemID,
		Price:       price,
		Description: description,
		ForSale:     forSale,
	}

	ctx := context.Background()

	if err := gorm.G[model.FestivalStock](r.db).Create(ctx, &stock); err != nil {
		return model.FestivalStock{}, wrapGormError(err)
	}

	return stock, nil
}

func (r *GormRepository) GetFestivalStockByID(festivalStockID uuid.UUID) (model.FestivalStock, error) {
	ctx := context.Background()

	stock, err := gorm.G[model.FestivalStock](r.db).
		Where(model.FestivalStock{ID: festivalStockID}, "ID").
		Preload("Festival", nil).
		Preload("StockItem", nil).
		First(ctx)

	if err != nil {
		return model.FestivalStock{}, wrapGormError(err)
	}

	return stock, nil
}

func (r *GormRepository) QueryFestivalStocks(festivalID uuid.UUID, category string, onlyForSale bool) ([]model.FestivalStock, error) {
	ctx := context.Background()

	query := gorm.G[model.FestivalStock](r.db).
		Joins(clause.JoinTarget{Association: "StockItem"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Where(model.StockItem{Category: category})
			return nil
		}).
		Joins(clause.JoinTarget{Association: "Festival"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Where(model.Festival{ID: festivalID})
			return nil
		})

	if onlyForSale {
		query = query.Where(model.FestivalStock{ForSale: true})
	}

	stocks, err := query.Find(ctx)

	if err != nil {
		return nil, wrapGormError(err)
	}

	return stocks, nil
}

func (r *GormRepository) UpdateFestivalStock(festivalStockID uuid.UUID, description string, forSale bool) error {
	ctx := context.Background()

	rows, err := gorm.G[model.FestivalStock](r.db).
		Where(model.FestivalStock{ID: festivalStockID}, "ID").
		Select("Description", "ForSale").
		Updates(ctx, model.FestivalStock{Description: description, ForSale: forSale})
	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *GormRepository) DeleteFestivalStock(festivalStockID uuid.UUID) error {
	ctx := context.Background()

	rows, err := gorm.G[model.FestivalStock](r.db).
		Where(model.FestivalStock{ID: festivalStockID}, "ID").
		Delete(ctx)
	if err != nil {
		return wrapGormError(err)
	}

	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}
