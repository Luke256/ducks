package repository

import (
	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type SaleRepository interface {
	// CreateSaleRecord 販売記録を作成します
	CreateSaleRecord(festivalStockID uuid.UUID, quantity int) (model.SaleRecord, error)

	// GetSaleRecordByID 販売記録IDから販売記録を取得します
	GetSaleRecordByID(saleRecordID uuid.UUID) (model.SaleRecord, error)

	// GetSaleRecordsByFestivalStockID イベント在庫IDから販売記録を取得します
	GetSaleRecordsByFestivalStockID(festivalStockID uuid.UUID) ([]model.SaleRecord, error)

	// QuerySaleRecords イベントIDと商品IDから販売記録を取得します
	QuerySaleRecords(festivalID, stockItemID uuid.UUID) ([]model.SaleRecord, error)

	// DeleteSaleRecord 販売記録を削除します
	DeleteSaleRecord(saleRecordID uuid.UUID) error
}