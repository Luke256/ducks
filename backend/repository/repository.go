package repository

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

type Repository interface {
	PosterRepository
	FestivalRepository
	StockItemRepository
}
