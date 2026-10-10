package gorm

import (
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterFestivalStock(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")

	for _, forSale := range []bool{true, false} {
		name := "Register Festival Stock For Sale"
		if !forSale {
			name = "Register Festival Stock Not For Sale"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			festivalStock, err := repo.RegisterFestivalStock(fes.ID, item.ID, 500, "Stock Description", forSale)
			require.NoError(t, err)
			assert.NotZero(t, festivalStock.ID)
			assert.Equal(t, fes.ID, festivalStock.FestivalID)
			assert.Equal(t, item.ID, festivalStock.StockItemID)
			assert.Equal(t, 500, festivalStock.Price)
			assert.Equal(t, "Stock Description", festivalStock.Description)
			assert.Equal(t, forSale, festivalStock.ForSale)

			retrievedStock, err := repo.GetFestivalStockByID(festivalStock.ID)
			require.NoError(t, err)
			assert.Equal(t, forSale, retrievedStock.ForSale)
		})
	}
}

func TestGetFestivalStockByID(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500, "Stock Description", true)

	t.Run("Get Festival Stock By ID", func(t *testing.T) {
		t.Parallel()
		retrievedStock, err := repo.GetFestivalStockByID(fesStock.ID)
		assert.NoError(t, err)
		assert.Equal(t, fesStock.ID, retrievedStock.ID)
		assert.Equal(t, fesStock.Price, retrievedStock.Price)
		assert.Equal(t, "Stock Description", retrievedStock.Description)
		assert.Equal(t, fesStock.ForSale, retrievedStock.ForSale)
		assert.Equal(t, fes.ID, retrievedStock.Festival.ID)
		assert.Equal(t, item.ID, retrievedStock.StockItem.ID)
	})

	t.Run("Get Non-Existent Festival Stock By ID", func(t *testing.T) {
		t.Parallel()
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		_, err = repo.GetFestivalStockByID(id)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Get Festival Stock By Zero UUID", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetFestivalStockByID(uuid.Nil)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestQueryFestivalStocks(t *testing.T) {
	t.Parallel()
	repo := setup(t, s2)

	fes1 := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	fes2 := mustCreateFestival(t, repo, "Another Fest for Stock", "Another Festival Description")
	item1 := mustCreateStockItem(t, repo, "Stock Item 1", "Item Description 1", "Category1", "image_id_1")
	item2 := mustCreateStockItem(t, repo, "Stock Item 2", "Item Description 2", "Category1", "image_id_2")
	item3 := mustCreateStockItem(t, repo, "Stock Item 3", "Item Description 3", "Category2", "image_id_3")
	stock1 := mustCreateFestivalStock(t, repo, fes1.ID, item1.ID, 500, "Stock Description 1", true)
	stock2 := mustCreateFestivalStock(t, repo, fes1.ID, item2.ID, 800, "Stock Description 2", false)
	stock3 := mustCreateFestivalStock(t, repo, fes1.ID, item3.ID, 1200, "Stock Description 3", true)
	stock4 := mustCreateFestivalStock(t, repo, fes2.ID, item1.ID, 700, "Stock Description 4", true)

	t.Run("Query All Festival Stocks", func(t *testing.T) {
		t.Parallel()
		stocks, err := repo.QueryFestivalStocks(uuid.Nil, "", false)
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
		assert.Equal(t, item1.ID, fetchedStock1.StockItem.ID)
		assert.Equal(t, 500, fetchedStock1.Price)
		assert.True(t, fetchedStock1.ForSale)
	})

	t.Run("Query Festival Stocks by Festival ID", func(t *testing.T) {
		t.Parallel()
		stocks, err := repo.QueryFestivalStocks(fes1.ID, "", false)
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
		t.Parallel()
		stocks, err := repo.QueryFestivalStocks(uuid.Nil, "Category1", false)
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
		t.Parallel()
		stocks, err := repo.QueryFestivalStocks(fes1.ID, "Category1", false)
		assert.NoError(t, err)
		assert.Len(t, stocks, 2)

		stockIDs := make(map[uuid.UUID]bool)
		for _, stock := range stocks {
			stockIDs[stock.ID] = true
		}
		assert.Contains(t, stockIDs, stock1.ID)
		assert.Contains(t, stockIDs, stock2.ID)
	})

	for _, tc := range []struct {
		name       string
		festivalID uuid.UUID
		category   string
		wantIDs    []uuid.UUID
	}{
		{"Query Only Stocks For Sale", uuid.Nil, "", []uuid.UUID{stock1.ID, stock3.ID, stock4.ID}},
		{"Query Stocks For Sale by Festival ID", fes1.ID, "", []uuid.UUID{stock1.ID, stock3.ID}},
		{"Query Stocks For Sale by Category", uuid.Nil, "Category1", []uuid.UUID{stock1.ID, stock4.ID}},
		{"Query Stocks For Sale by Festival ID and Category", fes1.ID, "Category1", []uuid.UUID{stock1.ID}},
		{"Query Stocks For Sale with No Matches", fes2.ID, "Category2", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stocks, err := repo.QueryFestivalStocks(tc.festivalID, tc.category, true)
			require.NoError(t, err)
			stockIDs := make([]uuid.UUID, 0, len(stocks))
			for _, stock := range stocks {
				assert.True(t, stock.ForSale)
				stockIDs = append(stockIDs, stock.ID)
			}
			assert.ElementsMatch(t, tc.wantIDs, stockIDs)
		})
	}
}

func TestUpdateFestivalStock(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	for _, forSale := range []bool{true, false} {
		name := "Update Festival Stock For Sale"
		if !forSale {
			name = "Update Festival Stock Not For Sale"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500, "Stock Description", !forSale)
			err := repo.UpdateFestivalStock(fesStock.ID, "Updated Stock Description", forSale)
			require.NoError(t, err)

			updatedStock, err := repo.GetFestivalStockByID(fesStock.ID)
			require.NoError(t, err)
			assert.Equal(t, 500, updatedStock.Price)
			assert.Equal(t, "Updated Stock Description", updatedStock.Description)
			assert.Equal(t, forSale, updatedStock.ForSale)
			assert.NoError(t, repo.UpdateFestivalStock(fesStock.ID, updatedStock.Description, updatedStock.ForSale))
		})
	}

	t.Run("Update Non-Existent Festival Stock", func(t *testing.T) {
		t.Parallel()
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		err = repo.UpdateFestivalStock(id, "Updated Stock Description", true)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Update Festival Stock with Zero UUID", func(t *testing.T) {
		t.Parallel()
		err := repo.UpdateFestivalStock(uuid.Nil, "Updated Stock Description", true)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteFestivalStock(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	fes := mustCreateFestival(t, repo, "Fest for Stock", "Festival Description")
	item := mustCreateStockItem(t, repo, "Stock Item", "Item Description", "Category", "image_id")
	fesStock := mustCreateFestivalStock(t, repo, fes.ID, item.ID, 500, "Stock Description", true)

	t.Run("Delete Festival Stock", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteFestivalStock(fesStock.ID)
		assert.NoError(t, err)

		_, err = repo.GetFestivalStockByID(fesStock.ID)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Festival Stock", func(t *testing.T) {
		t.Parallel()
		id, err := uuid.NewV7()
		assert.NoError(t, err)
		err = repo.DeleteFestivalStock(id)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Festival Stock with Zero UUID", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteFestivalStock(uuid.Nil)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
