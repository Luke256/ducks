package stockitem

import (
	"fmt"
	"mime/multipart"

	"github.com/Luke256/ducks/model"
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

func (m *ManagerImpl) toStockItemModel(item model.StockItem) StockItem {
	return StockItem{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Category:    item.Category,
		ImageURL:    m.storage.GetFileURL(item.ImageID),
	}
}

func (m *ManagerImpl) Create(name string, description string, category string, image *multipart.FileHeader) (_ StockItem, err error) {
	imageID, err := m.storage.UploadFile(image)
	if err != nil {
		return StockItem{}, fmt.Errorf("failed to upload image: %w", err)
	}
	defer func() {
		if err != nil {
			_ = m.storage.DeleteFile(imageID)
		}
	}()

	item, err := m.repo.RegisterStockItem(name, description, category, imageID)
	if err != nil {
		return StockItem{}, fmt.Errorf("failed to register stock item: %w", err)
	}

	return m.toStockItemModel(item), nil
}

func (m *ManagerImpl) Get(id uuid.UUID) (StockItem, error) {
	item, err := m.repo.GetStockItemByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return StockItem{}, ErrNotFound
		default:
			return StockItem{}, fmt.Errorf("failed to get stock item: %w", err)
		}
	}

	return m.toStockItemModel(item), nil
}

func (m *ManagerImpl) Query(category string) ([]StockItem, error) {
	items, err := m.repo.QueryStockItems(category)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock items: %w", err)
	}

	var result []StockItem
	for _, item := range items {
		result = append(result, m.toStockItemModel(item))
	}

	return result, nil
}

func (m *ManagerImpl) Edit(id uuid.UUID, name string, description string, category string) (err error) {
	item, err := m.repo.GetStockItemByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to get stock item: %w", err)
		}
	}

	_, err = m.repo.UpdateStockItem(id, name, description, category, item.ImageID)
	if err != nil {
		return fmt.Errorf("failed to update stock item: %w", err)
	}

	return nil
}

func (m *ManagerImpl) UpdateImage(id uuid.UUID, image *multipart.FileHeader) (err error) {
	item, err := m.repo.GetStockItemByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to get stock item: %w", err)
		}
	}

	if err := m.storage.DeleteFile(item.ImageID); err != nil {
		return fmt.Errorf("failed to delete old image from storage: %w", err)
	}

	imageID, err := m.storage.UploadFile(image)
	if err != nil {
		return fmt.Errorf("failed to upload new image: %w", err)
	}
	defer func() {
		if err != nil {
			_ = m.storage.DeleteFile(imageID)
		}
	}()

	_, err = m.repo.UpdateStockItem(id, item.Name, item.Description, item.Category, imageID)
	if err != nil {
		return fmt.Errorf("failed to update stock item image: %w", err)
	}

	return nil
}

func (m *ManagerImpl) Delete(id uuid.UUID) error {
	item, err := m.repo.GetStockItemByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to get stock item: %w", err)
		}
	}

	err = m.storage.DeleteFile(item.ImageID)
	if err != nil {
		return fmt.Errorf("failed to delete image from storage: %w", err)
	}

	err = m.repo.DeleteStockItem(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return fmt.Errorf("failed to delete stock item: %w", err)
		}
	}

	return nil
}