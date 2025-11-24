package festivalstock

import (
	"errors"

	stockitem "github.com/Luke256/ducks/service/stock_item"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("not found")
)

type Stock struct {
	ID          uuid.UUID           `json:"id"`
	Item        stockitem.StockItem `json:"item"`
	FestivalID  uuid.UUID           `json:"festival_id"`
	Price       int                 `json:"price"`
	Description string              `json:"description"`
}

type Manager interface {
	// Create イベントで販売するアイテムを登録します
	Create(festivalID, itemID uuid.UUID, price int, description string) (Stock, error)

	// Get 指定されたIDのイベントで販売するアイテムを取得します
	Get(id uuid.UUID) (Stock, error)

	// Query イベントIDやカテゴリで販売するアイテムを検索します
	// festivalID, categoryが空文字の場合、全てのカテゴリを対象とします
	Query(festivalID uuid.UUID, category string) ([]Stock, error)

	// Update 指定されたIDのイベントで販売するアイテムの説明を更新します
	Update(id uuid.UUID, description string) error

	// Delete 指定されたIDのイベントで販売するアイテムを削除します
	Delete(id uuid.UUID) error
}
