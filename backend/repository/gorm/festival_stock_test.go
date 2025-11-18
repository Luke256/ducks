package gorm

import (
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterFestivalStock(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")

	t.Run("Register Festival Stock", func(t *testing.T) {
		festivalStock, err := repo.RegisterFestivalStock(fes.ID, item.ID, 500)
		assert.NoError(t, err)
		assert.NotZero(t, festivalStock.ID)
		assert.Equal(t, fes.ID, festivalStock.FestivalID)
		assert.Equal(t, item.ID, festivalStock.StockItemID)
		assert.Equal(t, 500, festivalStock.Price)
	})
}

func TestGetFestivalStockByID(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500)

	t.Run("Get Festival Stock By ID", func(t *testing.T) {
		retrievedStock, err := repo.GetFestivalStockByID(fesStock.ID)
		assert.NoError(t, err)
		assert.Equal(t, fesStock.ID, retrievedStock.ID)
		assert.Equal(t, fesStock.Price, retrievedStock.Price)
		assert.Equal(t, fes.ID, retrievedStock.Festival.ID)
		assert.Equal(t, item.ID, retrievedStock.Stock.ID)
	})

	t.Run("Get Non-Existent Festival Stock By ID", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		_, err = repo.GetFestivalStockByID(id)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Get Festival Stock By Zero UUID", func(t *testing.T) {
		_, err := repo.GetFestivalStockByID(uuid.Nil)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestQueryFestivalStocks(t *testing.T) {
	repo := setup(t, s1)

	fes1 := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	fes2 := mustCreateFestival(t, repo, "Another Fest for Stock", "Another Festival Description")
	item1 := mustCreateStockItem(t, repo, "Stock Item 1", "Item Description 1", "Category1", "image_id_1")
	item2 := mustCreateStockItem(t, repo, "Stock Item 2", "Item Description 2", "Category1", "image_id_2")
	item3 := mustCreateStockItem(t, repo, "Stock Item 3", "Item Description 3", "Category2", "image_id_3")
	stock1 := mustCreateFestivalStock(t, repo, fes1.ID, item1.ID, 500)
	stock2 := mustCreateFestivalStock(t, repo, fes1.ID, item2.ID, 800)
	stock3 := mustCreateFestivalStock(t, repo, fes1.ID, item3.ID, 1200)
	stock4 := mustCreateFestivalStock(t, repo, fes2.ID, item1.ID, 700)

	t.Run("Query All Festival Stocks", func(t *testing.T) {
		stocks, err := repo.QueryFestivalStocks(uuid.Nil, "")
		assert.NoError(t, err)
		assert.Len(t, stocks, 4)

		stockIDs := make(map[uuid.UUID]bool)
		var fetchedStock1 model.FestivalStock
		for _, stock := range stocks {
			stockIDs[stock.ID] = true
			if stock.ID == stock1.ID {
				fetchedStock1 = stock
			}
		}
		assert.Contains(t, stockIDs, stock1.ID)
		assert.Contains(t, stockIDs, stock2.ID)
		assert.Contains(t, stockIDs, stock3.ID)
		assert.Contains(t, stockIDs, stock4.ID)

		assert.Equal(t, fes1.ID, fetchedStock1.Festival.ID)
		assert.Equal(t, item1.ID, fetchedStock1.Stock.ID)
		assert.Equal(t, 500, fetchedStock1.Price)
	})

	t.Run("Query Festival Stocks by Festival ID", func(t *testing.T) {
		stocks, err := repo.QueryFestivalStocks(fes1.ID, "")
		assert.NoError(t, err)
		assert.Len(t, stocks, 3)

		stockIDs := make(map[uuid.UUID]bool)
		for _, stock := range stocks {
			stockIDs[stock.ID] = true
		}
		assert.Contains(t, stockIDs, stock1.ID)
		assert.Contains(t, stockIDs, stock2.ID)
		assert.Contains(t, stockIDs, stock3.ID)
	})

	t.Run("Query Festival Stocks by Category", func(t *testing.T) {
		stocks, err := repo.QueryFestivalStocks(uuid.Nil, "Category1")
		assert.NoError(t, err)
		assert.Len(t, stocks, 3)

		stockIDs := make(map[uuid.UUID]bool)
		for _, stock := range stocks {
			stockIDs[stock.ID] = true
		}
		assert.Contains(t, stockIDs, stock1.ID)
		assert.Contains(t, stockIDs, stock2.ID)
		assert.Contains(t, stockIDs, stock4.ID)
	})

	t.Run("Query Festival Stocks by Festival ID and Category", func(t *testing.T) {
		stocks, err := repo.QueryFestivalStocks(fes1.ID, "Category1")
		assert.NoError(t, err)
		assert.Len(t, stocks, 2)

		stockIDs := make(map[uuid.UUID]bool)
		for _, stock := range stocks {
			stockIDs[stock.ID] = true
		}
		assert.Contains(t, stockIDs, stock1.ID)
		assert.Contains(t, stockIDs, stock2.ID)
	})
}

func TestUpdateFestivalStockPrice(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500)

	t.Run("Update Festival Stock Price", func(t *testing.T) {
		err := repo.UpdateFestivalStockPrice(fesStock.ID, 750)
		assert.NoError(t, err)

		updatedStock, err := repo.GetFestivalStockByID(fesStock.ID)
		assert.NoError(t, err)
		assert.Equal(t, 750, updatedStock.Price)
	})

	t.Run("Update Non-Existent Festival Stock Price", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		err = repo.UpdateFestivalStockPrice(id, 900)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Update Festival Stock Price with Zero UUID", func(t *testing.T) {
		err := repo.UpdateFestivalStockPrice(uuid.Nil, 900)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteFestivalStock(t *testing.T) {
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500)

	t.Run("Delete Festival Stock", func(t *testing.T) {
		err := repo.DeleteFestivalStock(fesStock.ID)
		assert.NoError(t, err)

		_, err = repo.GetFestivalStockByID(fesStock.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Festival Stock", func(t *testing.T) {
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		err = repo.DeleteFestivalStock(id)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Festival Stock with Zero UUID", func(t *testing.T) {
		err := repo.DeleteFestivalStock(uuid.Nil)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
