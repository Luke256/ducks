package repository

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrForeignKey    = errors.New("foreign key constraint failed")
)

type Repository interface {
	PosterRepository
	FestivalRepository
	StockItemRepository
	FestivalStockRepository
	SaleRepository
}
