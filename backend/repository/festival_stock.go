package repository

import (
	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type FestivalStockRepository interface {
	// RegisterFestivalStock イベントで販売するアイテムを登録します
	RegisterFestivalStock(festivalID, itemID uuid.UUID, price int) (model.FestivalStock, error)

	// GetFestivalStockByID イベントで販売するアイテムをIDで取得します
	GetFestivalStockByID(festivalStockID uuid.UUID) (model.FestivalStock, error)

	// QueryFestivalStocks イベントIDやカテゴリで販売するアイテムを検索します
	QueryFestivalStocks(festivalID uuid.UUID, category string) ([]model.FestivalStock, error)

	// UpdateFestivalStockPrice イベントで販売するアイテムの価格を更新します
	UpdateFestivalStockPrice(festivalStockID uuid.UUID, newPrice int) error

	// DeleteFestivalStock イベントで販売するアイテムを削除します
	DeleteFestivalStock(festivalStockID uuid.UUID) error
}