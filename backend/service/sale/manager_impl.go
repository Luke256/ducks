package sale

import (
	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	festivalstock "github.com/Luke256/ducks/service/festival_stock"

	"github.com/google/uuid"
)

type ManagerImpl struct {
	repo repository.Repository
}

func NewManagerImpl(saleRepo repository.Repository) *ManagerImpl {
	return &ManagerImpl{repo: saleRepo}
}

func (m *ManagerImpl) toSaleRecordType(record model.SaleRecord) SaleRecord {
	return SaleRecord{
		ID:       record.ID,
		StockID:  record.FestivalStockID,
		Quantity: record.Quantity,
	}
}

func (m *ManagerImpl) Create(stockID uuid.UUID, quantity int) (SaleRecord, error) {
	_, err := m.repo.GetFestivalStockByID(stockID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return SaleRecord{}, festivalstock.ErrNotFound
		default:
			return SaleRecord{}, err
		}
	}

	record, err := m.repo.CreateSaleRecord(stockID, quantity)
	if err != nil {
		return SaleRecord{}, err
	}
	return m.toSaleRecordType(record), nil
}

func (m *ManagerImpl) Get(id uuid.UUID) (SaleRecord, error) {
	record, err := m.repo.GetSaleRecordByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return SaleRecord{}, ErrNotFound
		default:
			return SaleRecord{}, err
		}
	}
	return m.toSaleRecordType(record), nil
}

func (m *ManagerImpl) GetByStockID(stockID uuid.UUID) ([]SaleRecord, error) {
	records, err := m.repo.GetSaleRecordsByFestivalStockID(stockID)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, festivalstock.ErrNotFound
		default:
			return nil, err
		}
	}

	result := make([]SaleRecord, len(records))
	for i, record := range records {
		result[i] = m.toSaleRecordType(record)
	}
	return result, nil
}

func (m *ManagerImpl) Query(festivalID, stockItemID uuid.UUID) ([]SaleRecord, error) {
	records, err := m.repo.QuerySaleRecords(festivalID, stockItemID)
	if err != nil {
		return nil, err
	}

	result := make([]SaleRecord, len(records))
	for i, record := range records {
		result[i] = m.toSaleRecordType(record)
	}
	return result, nil
}

func (m *ManagerImpl) Delete(id uuid.UUID) error {
	err := m.repo.DeleteSaleRecord(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}