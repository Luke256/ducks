package gorm

import (
	"github.com/Luke256/ducks/repository"
	"gorm.io/gorm"
)

func wrapGormError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return repository.ErrNotFound
	default:
		return err
	}
}
