package repository

import (
	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type FestivalRepository interface {
	// RegisterFestival イベントを登録します
	RegisterFestival(name string, description string) (model.Festival, error)

	// GetFestivalByID イベントIDからイベントを取得します
	GetFestivalByID(festivalID uuid.UUID) (model.Festival, error)

	// GetAllFestivals すべてのイベントを取得します
	GetAllFestivals() ([]model.Festival, error)

	// UpdateFestival イベント情報を更新します
	UpdateFestival(festivalID uuid.UUID, name string, description string) error

	// SetFestivalLocked イベントの変更禁止状態を設定・解除します
	SetFestivalLocked(festivalID uuid.UUID, locked bool) error

	// DeleteFestival イベントを削除します
	DeleteFestival(festivalID uuid.UUID) error
}
