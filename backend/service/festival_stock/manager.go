package festivalstock

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("not found")
)

type Stock struct {
	ID          string    `json:"id"`
	StockItemID uuid.UUID `json:"stock_item_id"`
	FestivalID  uuid.UUID `json:"festival_id"`
	Price       int       `json:"price"`
}

type Manager interface {
	// Create イベントで販売するアイテムを登録します
	Create(festivalID, itemID uuid.UUID, price int) (Stock, error)

	// Get 指定されたIDのイベントで販売するアイテムを取得します
	Get(id uuid.UUID) (Stock, error)

	// Query イベントIDやカテゴリで販売するアイテムを検索します
	// festivalID, categoryが空文字の場合、全てのカテゴリを対象とします
	Query(festivalID uuid.UUID, category string) ([]Stock, error)

	// UpdatePrice 指定されたIDのイベントで販売するアイテムの価格を更新します
	UpdatePrice(id uuid.UUID, newPrice int) error

	// Delete 指定されたIDのイベントで販売するアイテムを削除します
	Delete(id uuid.UUID) error
}
