package festival

import (
	"fmt"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
)

type FestivalImpl struct {
	repo repository.FestivalRepository
}

func NewFestivalManager(repo repository.FestivalRepository) *FestivalImpl {
	return &FestivalImpl{
		repo: repo,
	}
}

func (f *FestivalImpl) Create(name, description string) (Festival, error) {
	festival, err := f.repo.RegisterFestival(name, description)
	if err != nil {
		return Festival{}, err
	}

	return Festival{
		ID:          festival.ID,
		Name:        festival.Name,
		Description: festival.Description,
	}, nil
}

func (f *FestivalImpl) Get(id uuid.UUID) (Festival, error) {
	festival, err := f.repo.GetFestivalByID(id)
	if err != nil {
		return Festival{}, err
	}

	return Festival{
		ID:          festival.ID,
		Name:        festival.Name,
		Description: festival.Description,
	}, nil
}

func (f *FestivalImpl) List() ([]Festival, error) {
	festivals, err := f.repo.GetAllFestivals()
	if err != nil {
		return nil, err
	}

	result := make([]Festival, len(festivals))

	for i, festival := range festivals {
		result[i] = Festival{
			ID:          festival.ID,
			Name:        festival.Name,
			Description: festival.Description,
		}
	}

	return result, nil
}

func (f *FestivalImpl) Edit(id uuid.UUID, name, description string) error {
	err := f.repo.UpdateFestival(id, name, description)

	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to update festival: %w", err)
		}
	}

	return nil
}

func (f *FestivalImpl) Delete(id uuid.UUID) error {
	err := f.repo.DeleteFestival(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to delete festival: %w", err)
		}
	}

	return nil
}