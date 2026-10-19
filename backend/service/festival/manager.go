package festival

import (
	"errors"
	"time"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrEventLocked = repository.ErrEventLocked
)

type Festival struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Locked      bool      `json:"locked"`
}

type Manager interface {
	// Create イベントを作成します
	Create(name, description string, date time.Time) (Festival, error)

	// Get 指定されたIDのイベントを取得します
	Get(id uuid.UUID) (Festival, error)

	// List 開催日の昇順でイベントを取得します（未設定は末尾）
	List() ([]Festival, error)

	// Edit 指定されたIDのイベント情報を更新します
	Edit(id uuid.UUID, name, description string, date time.Time) error

	// SetLocked イベントの変更禁止状態を設定・解除します
	SetLocked(id uuid.UUID, locked bool) error

	// Delete 指定されたIDのイベントを削除します
	Delete(id uuid.UUID) error
}
