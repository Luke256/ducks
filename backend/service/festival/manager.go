package festival

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound	  = errors.New("not found")
)

type Festival struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Manager interface {
	// Create イベントを作成します
	Create(name, description string) (Festival, error)

	// Get 指定されたIDのイベントを取得します
	Get(id uuid.UUID) (Festival, error)

	// List すべてのイベントを取得します
	List() ([]Festival, error)

	// Edit 指定されたIDのイベント情報を更新します
	Edit(id uuid.UUID, name, description string) error

	// Delete 指定されたIDのイベントを削除します
	Delete(id uuid.UUID) error
}
