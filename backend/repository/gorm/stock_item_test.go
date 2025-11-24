package gorm

import (
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterStockItem(t *testing.T) {
	repo := setup(t, common)

	t.Run("Register Stock Item", func(t *testing.T) {
		item, err := repo.RegisterStockItem("Test Item", "This is a test item", "Category1", "img-123")

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, item.ID)
		assert.Equal(t, "Test Item", item.Name)
		assert.Equal(t, "This is a test item", item.Description)
		assert.Equal(t, "Category1", item.Category)
		assert.Equal(t, "img-123", item.ImageID)
	})
}

func TestGetStockItemByID(t *testing.T) {
	repo := setup(t, common)

	item := mustCreateStockItem(t, repo, "Sample Item", "Sample Description", "CategoryA", "img-456")

	t.Run("Get Stock Item By ID", func(t *testing.T) {
		fetchedItem, err := repo.GetStockItemByID(item.ID)

		assert.NoError(t, err)
		assert.Equal(t, item.ID, fetchedItem.ID)
		assert.Equal(t, item.Name, fetchedItem.Name)
		assert.Equal(t, item.Description, fetchedItem.Description)
		assert.Equal(t, item.Category, fetchedItem.Category)
		assert.Equal(t, item.ImageID, fetchedItem.ImageID)
	})

	t.Run("Get Stock Item By Invalid ID", func(t *testing.T) {
		_, err := repo.GetStockItemByID(uuid.New())

		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Get Stock Item By Zero UUID", func(t *testing.T) {
		_, err := repo.GetStockItemByID(uuid.Nil)

		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestQueryStockItems(t *testing.T) {
	repo := setup(t, s1)

	var createdItems = []model.StockItem{
		mustCreateStockItem(t, repo, "Item 1", "Desc 1", "Cat1", "img-1"),
		mustCreateStockItem(t, repo, "Item 2", "Desc 2", "Cat1", "img-2"),
		mustCreateStockItem(t, repo, "Item 3", "Desc 3", "Cat2", "img-3"),
	}

	t.Run("Query All Items", func(t *testing.T) {
		items, err := repo.QueryStockItems("")
		
		assert.NoError(t, err)
		assert.Len(t, items, 3)
		
		itemIDs := make(map[uuid.UUID]bool)
		for _, item := range items {
			itemIDs[item.ID] = true

			for _, createdItem := range createdItems {
				if item.ID == createdItem.ID {
					assert.Equal(t, createdItem.Name, item.Name)
					assert.Equal(t, createdItem.Description, item.Description)
					assert.Equal(t, createdItem.Category, item.Category)
					assert.Equal(t, createdItem.ImageID, item.ImageID)
				}
			}
		}
		
		for _, createdItem := range createdItems {
			assert.True(t, itemIDs[createdItem.ID], "Expected item ID not found: %s", createdItem.ID)
		}
	})

	t.Run("Query Items by Category", func(t *testing.T) {
		items, err := repo.QueryStockItems("Cat1")
		
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		
		itemIDs := make(map[uuid.UUID]bool)
		for _, item := range items {
			itemIDs[item.ID] = true
		}

		for _, createdItem := range createdItems[:2] {
			assert.True(t, itemIDs[createdItem.ID], "Expected item ID not found: %s", createdItem.ID)
		}
	})
}

func TestUpdateStockItem(t *testing.T) {
	repo := setup(t, common)

	item := mustCreateStockItem(t, repo, "Old Name", "Old Description", "OldCategory", "old-img")

	t.Run("Update Stock Item", func(t *testing.T) {
		updatedItem, err := repo.UpdateStockItem(item.ID, "New Name", "New Description", "NewCategory", "new-img")

		assert.NoError(t, err)
		assert.Equal(t, item.ID, updatedItem.ID)
		assert.Equal(t, "New Name", updatedItem.Name)
		assert.Equal(t, "New Description", updatedItem.Description)
		assert.Equal(t, "NewCategory", updatedItem.Category)
		assert.Equal(t, "new-img", updatedItem.ImageID)
	})

	t.Run("Update Non-Existent Stock Item", func(t *testing.T) {
		_, err := repo.UpdateStockItem(uuid.New(), "Name", "Description", "Category", "img")

		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Update Stock Item with Zero UUID", func(t *testing.T) {
		_, err := repo.UpdateStockItem(uuid.Nil, "Name", "Description", "Category", "img")

		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteStockItem(t *testing.T) {
	repo := setup(t, common)

	item := mustCreateStockItem(t, repo, "To Be Deleted", "Description", "Category", "img-del")

	t.Run("Delete Stock Item", func(t *testing.T) {
		err := repo.DeleteStockItem(item.ID)

		assert.NoError(t, err)

		_, err = repo.GetStockItemByID(item.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Stock Item", func(t *testing.T) {
		err := repo.DeleteStockItem(uuid.New())

		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Stock Item with Zero UUID", func(t *testing.T) {
		err := repo.DeleteStockItem(uuid.Nil)

		assert.Equal(t, repository.ErrNotFound, err)
	})
}