package sale

import (
	"errors"

	"github.com/google/uuid"
)

type SaleRecord struct {
	ID        uuid.UUID `json:"id"`
	StockID   uuid.UUID `json:"stock_id"`
	Quantity  int       `json:"quantity"`
}

var (
	ErrNotFound = errors.New("sale record not found")
)

type SaleManager interface {
	// Create 購入記録を作成します
	Create(StockID uuid.UUID, quantity int) (SaleRecord, error)

	// Get 購入記録をIDで取得します
	Get(id uuid.UUID) (SaleRecord, error)

	// GetByStockID 商品IDで購入記録を取得します
	GetByStockID(stockID uuid.UUID) ([]SaleRecord, error)

	// Query 購入記録を検索します
	Query(festivalID, stockItemID uuid.UUID) ([]SaleRecord, error)

	// Delete 購入記録を削除します
	Delete(id uuid.UUID) error
}