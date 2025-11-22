package gorm

import (
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

func (r *GormRepository) CreateSaleRecord(festivalStockID uuid.UUID, quantity int) (model.SaleRecord, error) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	if err != nil {
		return model.SaleRecord{}, err
	}
	saleRecord := model.SaleRecord{
		ID:              id,
		FestivalStockID: festivalStockID,
		Quantity:        quantity,
	}

	if err := gorm.G[model.SaleRecord](r.db).Create(ctx, &saleRecord); err != nil {
		return model.SaleRecord{}, wrapGormError(err)
	}

	return saleRecord, nil
}

func (r *GormRepository) GetSaleRecordByID(saleRecordID uuid.UUID) (model.SaleRecord, error) {
	ctx := context.Background()
	saleRecord, err := gorm.G[model.SaleRecord](r.db).
		Where(model.SaleRecord{ID: saleRecordID}, "ID").
		First(ctx)
	if err != nil {
		return model.SaleRecord{}, wrapGormError(err)
	}
	return saleRecord, nil
}

func (r *GormRepository) GetSaleRecordsByFestivalStockID(festivalStockID uuid.UUID) ([]model.SaleRecord, error) {
	ctx := context.Background()

	saleRecords, err := gorm.G[model.SaleRecord](r.db).
		Where(model.SaleRecord{FestivalStockID: festivalStockID}, "FestivalStockID").
		Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}

	return saleRecords, nil
}

func (r *GormRepository) QuerySaleRecords(festivalID, stockItemID uuid.UUID) ([]model.SaleRecord, error) {
	ctx := context.Background()

	saleRecords, err := gorm.G[model.SaleRecord](r.db).
		Joins(clause.JoinTarget{Association: "FestivalStock"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Where(model.FestivalStock{
				FestivalID:  festivalID,
				StockItemID: stockItemID,
			})
			return nil
		}).
		Find(ctx)
	if err != nil {
		return nil, wrapGormError(err)
	}

	return saleRecords, nil
}

func (r *GormRepository) DeleteSaleRecord(saleRecordID uuid.UUID) error {
	ctx := context.Background()
	rows, err := gorm.G[model.SaleRecord](r.db).
		Where(model.SaleRecord{ID: saleRecordID}, "ID").
		Delete(ctx)
	if err != nil {
		return wrapGormError(err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}
