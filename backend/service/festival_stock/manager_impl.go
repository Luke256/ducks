package festivalstock

import (
	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
)

type ManagerImpl struct {
	repo repository.FestivalStockRepository
}

func NewManager(repo repository.FestivalStockRepository) *ManagerImpl {
	return &ManagerImpl{repo: repo}
}

func (fm *ManagerImpl) toStockType(fs model.FestivalStock) Stock {
	return Stock{
		ID:          fs.ID.String(),
		StockItemID: fs.StockItemID,
		FestivalID:  fs.FestivalID,
		Price:       fs.Price,
	}
}

func (fm *ManagerImpl) Create(festivalID, itemID uuid.UUID, price int) (Stock, error) {
	fesStock, err := fm.repo.RegisterFestivalStock(festivalID, itemID, price)
	if err != nil {
		return Stock{}, err
	}

	return fm.toStockType(fesStock), nil
}

func (fm *ManagerImpl) Get(id uuid.UUID) (Stock, error) {
	fesStock, err := fm.repo.GetFestivalStockByID(id)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return Stock{}, ErrNotFound
		default:
			return Stock{}, err
		}
	}

	return fm.toStockType(fesStock), nil
}

func (fm *ManagerImpl) Query(festivalID uuid.UUID, category string) ([]Stock, error) {
	fesStocks, err := fm.repo.QueryFestivalStocks(festivalID, category)
	if err != nil {
		return nil, err
	}

	result := make([]Stock, len(fesStocks))
	for i, fs := range fesStocks {
		result[i] = fm.toStockType(fs)
	}

	return result, nil
}

func (fm *ManagerImpl) UpdatePrice(id uuid.UUID, newPrice int) error {
	err := fm.repo.UpdateFestivalStockPrice(id, newPrice)
	switch err {
	case nil:
		return nil
	case repository.ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (fm *ManagerImpl) Delete(id uuid.UUID) error {
	err := fm.repo.DeleteFestivalStock(id)
	switch err {
	case nil:
		return nil
	case repository.ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}