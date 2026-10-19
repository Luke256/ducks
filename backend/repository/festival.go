package repository

import (
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type FestivalRepository interface {
	// RegisterFestival イベントを登録します
	RegisterFestival(name string, description string, date time.Time) (model.Festival, error)

	// GetFestivalByID イベントIDからイベントを取得します
	GetFestivalByID(festivalID uuid.UUID) (model.Festival, error)

	// GetAllFestivals 開催日の昇順でイベントを取得します（未設定は末尾）
	GetAllFestivals() ([]model.Festival, error)

	// UpdateFestival イベント情報を更新します
	UpdateFestival(festivalID uuid.UUID, name string, description string, date time.Time) error

	// SetFestivalLocked イベントの変更禁止状態を設定・解除します
	SetFestivalLocked(festivalID uuid.UUID, locked bool) error

	// DeleteFestival イベントを削除します
	DeleteFestival(festivalID uuid.UUID) error
}
