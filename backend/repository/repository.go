package repository

import "errors"

var (
	ErrNotFound = errors.New("poster not found")
)

type Repository interface {
	ImageRepository
	FestivalRepository
}