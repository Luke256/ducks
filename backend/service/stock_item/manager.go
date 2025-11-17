package stockitem

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
)

type StockItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
}

var (
	ErrNotFound      = errors.New("not found")
)

type Manager interface {
	// Create アイテムを作成します
	Create(name string, description string, category string, image *multipart.FileHeader) (StockItem, error)

	// Get 指定されたIDのアイテムを取得します
	Get(id uuid.UUID) (StockItem, error)

	// Query アイテムを検索します
	// categoryが空文字の場合、全てのアイテムを取得します
	Query(category string) ([]StockItem, error)

	// Edit 指定されたIDのアイテム情報を更新します
	Edit(id uuid.UUID, name string, description string, category string) error

	// UpdateImage 指定されたIDのアイテムの画像を更新します
	UpdateImage(id uuid.UUID, image *multipart.FileHeader) error

	// Delete 指定されたIDのアイテムを削除します
	Delete(id uuid.UUID) error
}