package poster

import (
	"fmt"
	"mime/multipart"

	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/utils/storage"
	"github.com/google/uuid"
)

type ManagerImpl struct {
	repo    repository.Repository
	storage storage.Storage
}

func NewManagerImpl(repo repository.Repository, storage storage.Storage) *ManagerImpl {
	return &ManagerImpl{repo: repo, storage: storage}
}

func (m *ManagerImpl) Create(name string, festivalID uuid.UUID, description string, image *multipart.FileHeader) (Poster, error) {
	imageID, err := m.storage.UploadFile(image)
	if err != nil {
		return Poster{}, fmt.Errorf("failed to upload image: %w", err)
	}

	// duplicate check
	_, err = m.repo.GetPosterByFestivalIDAndPosterName(festivalID, name)
	if err == nil {
		return Poster{}, ErrAlreadyExists
	}
	if err != repository.ErrNotFound {
		return Poster{}, fmt.Errorf("failed to check duplicate poster: %w", err)
	}

	// festival existence check
	_, err = m.repo.GetFestivalByID(festivalID)
	if err != nil {
		if err == repository.ErrNotFound {
			return Poster{}, ErrNotFound
		}
		return Poster{}, fmt.Errorf("failed to check festival existence: %w", err)
	}

	poster, err := m.repo.RegisterPoster(festivalID, name, description, imageID)
	if err != nil {
		return Poster{}, fmt.Errorf("failed to register poster: %w", err)
	}

	return Poster{
		ID:          poster.ID,
		Name:        poster.PosterName,
		FestivalID:  poster.FestivalID,
		Description: poster.Description,
		ImageURL:    m.storage.GetFileURL(poster.ImageID),
		Status:      poster.Status,
	}, nil
}

func (m *ManagerImpl) Get(id uuid.UUID) (Poster, error) {
	poster, err := m.repo.GetPosterByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return Poster{}, ErrNotFound
		default:
			return Poster{}, fmt.Errorf("failed to get poster by ID: %w", err)
		}
	}

	return Poster{
		ID:          poster.ID,
		Name:        poster.PosterName,
		FestivalID:  poster.FestivalID,
		Description: poster.Description,
		ImageURL:    m.storage.GetFileURL(poster.ImageID),
		Status:      poster.Status,
	}, nil
}

func (m *ManagerImpl) GetByFestival(festivalID uuid.UUID) ([]Poster, error) {
	posters, err := m.repo.GetPostersByFestivalID(festivalID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get posters by festival ID: %w", err)
		}
	}

	result := make([]Poster, len(posters))
	for i, p := range posters {
		result[i] = Poster{
			ID:          p.ID,
			Name:        p.PosterName,
			FestivalID:  p.FestivalID,
			Description: p.Description,
			ImageURL:    m.storage.GetFileURL(p.ImageID),
			Status:      p.Status,
		}
	}

	return result, nil
}

func (m *ManagerImpl) GetByName(festivalID uuid.UUID, name string) (Poster, error) {
	poster, err := m.repo.GetPosterByFestivalIDAndPosterName(festivalID, name)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return Poster{}, ErrNotFound
		default:
			return Poster{}, fmt.Errorf("failed to get poster by festival ID and name: %w", err)
		}
	}

	return Poster{
		ID:          poster.ID,
		Name:        poster.PosterName,
		FestivalID:  poster.FestivalID,
		Description: poster.Description,
		ImageURL:    m.storage.GetFileURL(poster.ImageID),
		Status:      poster.Status,
	}, nil
}

func (m *ManagerImpl) Edit(id uuid.UUID, name, description string) error {
	err := m.repo.UpdatePoster(id, name, description)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to update poster: %w", err)
		}
	}
	return nil
}

func (m *ManagerImpl) ChangeStatus(id uuid.UUID, status string) error {
	err := m.repo.UpdatePosterStatus(id, status)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to update poster status: %w", err)
		}
	}
	return nil
}

func (m *ManagerImpl) Delete(id uuid.UUID) error {
	// delete poster image from storage
	poster, err := m.repo.GetPosterByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to get poster by ID: %w", err)
		}
	}
	
	err = m.storage.DeleteFile(poster.ImageID)
	if err != nil {
		return fmt.Errorf("failed to delete poster image from storage: %w", err)
	}

	err = m.repo.DeletePoster(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to delete poster: %w", err)
		}
	}
	return nil
}
