package poster

import (
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	"github.com/Luke256/ducks/utils/compressor"
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

func (m *ManagerImpl) Create(name string, festivalID uuid.UUID, description string, images []*multipart.FileHeader) (_ Poster, err error) {
	if len(images) == 0 {
		return Poster{}, ErrInvalidImages
	}
	if err := validateUploads(images); err != nil {
		return Poster{}, err
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
	fes, err := m.repo.GetFestivalByID(festivalID)
	if err != nil {
		if err == repository.ErrNotFound {
			return Poster{}, ErrNotFound
		}
		return Poster{}, fmt.Errorf("failed to check festival existence: %w", err)
	}

	imageIDs, err := m.uploadImages(images)
	if err != nil {
		return Poster{}, err
	}
	defer func() {
		if err != nil {
			m.deleteFiles(imageIDs)
		}
	}()
	poster, err := m.repo.RegisterPoster(festivalID, name, description, imageIDs)
	if err != nil {
		return Poster{}, fmt.Errorf("failed to register poster: %w", err)
	}

	return Poster{
		ID:          poster.ID,
		Name:        poster.PosterName,
		Description: poster.Description,
		Images:      m.posterImages(poster.Images),
		Status:      poster.Status,
		Festival:    festival.Festival{ID: fes.ID, Name: fes.Name, Description: fes.Description},
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
		Description: poster.Description,
		Images:      m.posterImages(poster.Images),
		Status:      poster.Status,
		Festival:    festival.Festival{ID: poster.Festival.ID, Name: poster.Festival.Name, Description: poster.Festival.Description},
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
			Description: p.Description,
			Images:      m.posterImages(p.Images),
			Status:      p.Status,
			Festival:    festival.Festival{ID: p.Festival.ID, Name: p.Festival.Name, Description: p.Festival.Description},
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
		Description: poster.Description,
		Images:      m.posterImages(poster.Images),
		Status:      poster.Status,
		Festival:    festival.Festival{ID: poster.Festival.ID, Name: poster.Festival.Name, Description: poster.Festival.Description},
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

func (m *ManagerImpl) UpdateImages(id uuid.UUID, images []*multipart.FileHeader, deleteIDs []string) (_ []Image, err error) {
	if len(images) == 0 && len(deleteIDs) == 0 {
		return nil, ErrInvalidImages
	}
	if err := validateUploads(images); err != nil {
		return nil, err
	}
	poster, err := m.repo.GetPosterByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get poster: %w", err)
	}
	owned := make(map[string]bool, len(poster.Images))
	for _, img := range poster.Images {
		owned[img.ID] = true
	}
	for _, imageID := range deleteIDs {
		if !owned[imageID] {
			return nil, ErrInvalidImages
		}
		delete(owned, imageID)
	}
	count := len(owned) + len(images)
	if count < 1 || count > repository.MaxPosterImages {
		return nil, ErrInvalidImages
	}

	imageIDs, err := m.uploadImages(images)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			m.deleteFiles(imageIDs)
		}
	}()
	updated, err := m.repo.UpdatePosterImages(id, imageIDs, deleteIDs)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to update poster images: %w", err)
	}
	// DBへの反映が確定してから、使わなくなったファイルを削除する。
	m.deleteFiles(deleteIDs)
	return m.posterImages(updated), nil
}

func (m *ManagerImpl) Delete(id uuid.UUID) error {
	poster, err := m.repo.GetPosterByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to get poster by ID: %w", err)
		}
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
	imageIDs := make([]string, len(poster.Images))
	for i, img := range poster.Images {
		imageIDs[i] = img.ID
	}
	m.deleteFiles(imageIDs)
	return nil
}

func (m *ManagerImpl) posterImages(images []model.PosterImage) []Image {
	result := make([]Image, len(images))
	for i, img := range images {
		result[i] = Image{ID: img.ID, URL: m.storage.GetFileURL(img.ID)}
	}
	return result
}

func validateUploads(images []*multipart.FileHeader) error {
	if len(images) > repository.MaxPosterImages {
		return ErrInvalidImages
	}
	for _, file := range images {
		if file == nil {
			return ErrInvalidImages
		}
		if file.Size > MaxImageSize {
			return ErrImageTooLarge
		}
	}
	return nil
}

func (m *ManagerImpl) uploadImages(images []*multipart.FileHeader) ([]string, error) {
	ids := make([]string, 0, len(images))
	for _, file := range images {
		id, err := m.storage.UploadFile(file)
		if err != nil {
			m.deleteFiles(ids)
			if errors.Is(err, compressor.ErrInvalidImage) {
				return nil, ErrInvalidImages
			}
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (m *ManagerImpl) deleteFiles(ids []string) {
	// shortcut: 削除失敗はログのみ。孤立ファイルが増える運用では再試行処理を追加する。
	for _, id := range ids {
		if err := m.storage.DeleteFile(id); err != nil {
			slog.Warn("failed to delete poster image", "image_id", id, "error", err)
		}
	}
}
