package gorm

import (
	"github.com/Luke256/ducks/model"

	"github.com/google/uuid"
)

func (r *GormRepository) RegisterFestivalStock(festivalID, itemID uuid.UUID, price int) (model.FestivalStock, error) {
	return model.FestivalStock{}, nil
}

func (r *GormRepository) GetFestivalStockByID(festivalStockID uuid.UUID) (model.FestivalStock, error) {
	return model.FestivalStock{}, nil
}

func (r *GormRepository) GetFestivalStocksByFestivalID(festivalID uuid.UUID) ([]model.FestivalStock, error) {
	return []model.FestivalStock{}, nil
}

func (r *GormRepository) UpdateFestivalStockPrice(festivalStockID uuid.UUID, newPrice int) error {
	return nil
}

func (r *GormRepository) DeleteFestivalStock(festivalStockID uuid.UUID) error {
	return nil
}

