package festival

import (
	"fmt"
	"time"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
)

type ManagerImpl struct {
	repo repository.FestivalRepository
}

func NewManagerImpl(repo repository.FestivalRepository) *ManagerImpl {
	return &ManagerImpl{
		repo: repo,
	}
}

func (f *ManagerImpl) Create(name, description string, date time.Time) (Festival, error) {
	festival, err := f.repo.RegisterFestival(name, description, date)
	if err != nil {
		return Festival{}, err
	}

	return Festival{
		ID:          festival.ID,
		Name:        festival.Name,
		Description: festival.Description,
		Locked:      festival.Locked,
		Date:        festival.Date,
	}, nil
}

func (f *ManagerImpl) Get(id uuid.UUID) (Festival, error) {
	festival, err := f.repo.GetFestivalByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return Festival{}, ErrNotFound
		default:
			return Festival{}, fmt.Errorf("failed to get festival: %w", err)
		}
	}

	return Festival{
		ID:          festival.ID,
		Name:        festival.Name,
		Description: festival.Description,
		Locked:      festival.Locked,
		Date:        festival.Date,
	}, nil
}

func (f *ManagerImpl) List() ([]Festival, error) {
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
			Locked:      festival.Locked,
			Date:        festival.Date,
		}
	}

	return result, nil
}

func (f *ManagerImpl) Edit(id uuid.UUID, name, description string, date time.Time) error {
	err := f.repo.UpdateFestival(id, name, description, date)

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

func (f *ManagerImpl) Delete(id uuid.UUID) error {
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

func (f *ManagerImpl) SetLocked(id uuid.UUID, locked bool) error {
	err := f.repo.SetFestivalLocked(id, locked)
	if err == repository.ErrNotFound {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to set festival lock: %w", err)
	}
	return nil
}
