package repository

import (
	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type StockItemRepository interface {
	// RegisterStockItem アイテムを登録します
	RegisterStockItem(name string, description string, category string, imageID string) (model.StockItem, error)

	// GetStockItemByID IDからアイテムを取得します
	GetStockItemByID(id uuid.UUID) (model.StockItem, error)

	// QueryStockItems アイテムを検索します
	QueryStockItems(category string) ([]model.StockItem, error)

	// UpdateStockItem アイテムを更新します
	UpdateStockItem(id uuid.UUID, name string, description string, category string, imageID string) (model.StockItem, error)

	// DeleteStockItem アイテムを削除します
	DeleteStockItem(id uuid.UUID) error
}
