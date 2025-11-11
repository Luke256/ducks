package gorm

import (
	"github.com/Luke256/ducks/repository"
	"gorm.io/gorm"
)

func wrapGormError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return repository.ErrPosterNotFound
	}

	return err
}