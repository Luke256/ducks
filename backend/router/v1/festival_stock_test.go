package v1

import (
	"testing"

	"github.com/google/uuid"
)

func TestRegisterFestivalStock(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	festival := env.mustCreateFestival(t, "Test Festival", "A festival for testing")
	stockItem := env.mustCreateStockItem(t, "Test Stock Item", "A stock item for testing", "Category1")

	t.Run("Register Festival Stock", func(t *testing.T) {
		res := e.POST("/api/festivals/{festival_id}/stocks", festival.ID).
			WithJSON(map[string]any{
				"item_id": stockItem.ID,
				"price":   1500,
				"description": "Stock Description",
			}).
			Expect().
			Status(201).
			JSON().
			Object()
			
		res.Value("festival_id").IsEqual(festival.ID.String())
		res.Value("price").IsEqual(1500)
		res.Value("description").IsEqual("Stock Description")
	})

	// Invalid Requests
	tests := []struct {
		name       string
		festivalID string
		payload    map[string]any
		expectCode int
	}{
		{
			name:       "Invalid Festival ID Format",
			festivalID: "invalid-uuid",
			payload: map[string]any{
				"item_id": stockItem.ID,
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name: 	 "Non-existent Festival ID",
			festivalID: uuid.New().String(),
			payload: map[string]any{
				"item_id": stockItem.ID,
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name:       "Zero Festival ID",
			festivalID: uuid.Nil.String(),
			payload: map[string]any{
				"item_id": stockItem.ID,
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name:       "Missing Item ID",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"price": 1500,
				"description": "Stock Description",
			},
			expectCode: 400,
		},
		{
			name:       "Invalid Item ID Format",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"item_id": "invalid-uuid",
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name:       "Non-existent Item ID",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"item_id": uuid.New().String(),
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name:       "Zero Item ID",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"item_id": uuid.Nil.String(),
				"price":   1500,
				"description": "Stock Description",
			},
			expectCode: 404,
		},
		{
			name:       "Missing Price",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"item_id": stockItem.ID,
				"description": "Stock Description",
			},
			expectCode: 400,
		},
		{
			name: "Empty Description",
			festivalID: festival.ID.String(),
			payload: map[string]any{
				"item_id": stockItem.ID,
				"price":   1500,
				"description": "",
			},
			expectCode: 201,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/api/festivals/{festival_id}/stocks", tc.festivalID).
				WithJSON(tc.payload).
				Expect().
				Status(tc.expectCode)
		})
	}
}

func TestGetFestivalStock(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "A festival for testing")
	item := env.mustCreateStockItem(t, "Test Stock Item", "A stock item for testing", "Category1")
	fesStock := env.mustCreateFestivalStock(t, fes.ID, item.ID, 2000, "Stock Description")

	t.Run("Get Festival Stock", func (t *testing.T) {
		res := e.GET("/api/stocks/{festival_stock_id}", fesStock.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("id").IsEqual(fesStock.ID.String())
		res.Value("festival_id").IsEqual(fes.ID.String())
		res.Value("item").Object().Value("id").IsEqual(item.ID.String())
		res.Value("price").IsEqual(2000)
		res.Value("description").IsEqual("Stock Description")
	})

	t.Run("Get Festival Stock - Not Found", func (t *testing.T) {
		e.GET("/api/stocks/{festival_stock_id}", uuid.New()).
			Expect().
			Status(404)
	})

	t.Run("Get Festival Stock - Invalid ID", func (t *testing.T) {
		e.GET("/api/stocks/{festival_stock_id}", "invalid-uuid").
			Expect().
			Status(404)
	})
}

func TestQueryFestivalStocks(t *testing.T) {
	env := setup(t, s2)
	e := env.R(t)

	fes1 := env.mustCreateFestival(t, "Festival One", "First festival")
	fes2 := env.mustCreateFestival(t, "Festival Two", "Second festival")

	itemA := env.mustCreateStockItem(t, "Item A", "First item", "Cat1")
	itemB := env.mustCreateStockItem(t, "Item B", "Second item", "Cat2")

	st1 := env.mustCreateFestivalStock(t, fes1.ID, itemA.ID, 1000, "Stock Description")
	st2 :=env.mustCreateFestivalStock(t, fes1.ID, itemB.ID, 1500, "Stock Description")
	env.mustCreateFestivalStock(t, fes2.ID, itemA.ID, 2000, "Stock Description")

	stock1, err := env.FSM.Get(st1.ID)
	if err != nil {
		t.Fatalf("Failed to get stock1: %v", err)
	}
	stock2, err := env.FSM.Get(st2.ID)
	if err != nil {
		t.Fatalf("Failed to get stock2: %v", err)
	}

	t.Run("Query Festival Stocks by Festival ID", func (t *testing.T) {
		res := e.GET("/api/festivals/{festival_id}/stocks", fes1.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("stocks").Array().Length().IsEqual(2)
		res.Value("stocks").Array().ContainsOnly(stock1, stock2)
	})

	t.Run("Query Festival Stocks - With Category Filter", func (t *testing.T) {
		res := e.GET("/api/festivals/{festival_id}/stocks", fes1.ID).
			WithQuery("category", "Cat1").
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("stocks").Array().Length().IsEqual(1)
		res.Value("stocks").Array().ContainsOnly(stock1)
	})
}

func TestUpdateFestivalStock(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "A festival for testing")
	item := env.mustCreateStockItem(t, "Test Stock Item", "A stock item for testing", "Category1")
	fesStock := env.mustCreateFestivalStock(t, fes.ID, item.ID, 2000, "Stock Description")

	t.Run("Update Festival Stock Price", func(t *testing.T) {
		e.PUT("/api/stocks/{festival_stock_id}", fesStock.ID).
			WithJSON(map[string]any{
				"description": "Updated Stock Description",
			}).
			Expect().
			Status(204)

		res := e.GET("/api/stocks/{festival_stock_id}", fesStock.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("description").IsEqual("Updated Stock Description")
	})

	t.Run("Update Festival Stock Price - Not Found", func(t *testing.T) {
		e.PUT("/api/stocks/{festival_stock_id}", uuid.New()).
			WithJSON(map[string]any{
				"description": "Updated Stock Description",
			}).
			Expect().
			Status(404)
	})

	t.Run("Update Festival Stock Price - Invalid ID", func(t *testing.T) {
		e.PUT("/api/stocks/{festival_stock_id}", "invalid-uuid").
			WithJSON(map[string]any{
				"description": "Updated Stock Description",
			}).
			Expect().
			Status(404)
	})

	t.Run("Update Festival Stock Price - Zero ID", func(t *testing.T) {
		e.PUT("/api/stocks/{festival_stock_id}", uuid.Nil).
			WithJSON(map[string]any{
				"description": "Updated Stock Description",
			}).
			Expect().
			Status(404)
	})
}

func TestDeleteFestivalStock(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "A festival for testing")
	item := env.mustCreateStockItem(t, "Test Stock Item", "A stock item for testing", "Category1")
	fesStock := env.mustCreateFestivalStock(t, fes.ID, item.ID, 2000, "Stock Description")

	t.Run("Delete Festival Stock - Not Implemented", func(t *testing.T) {
		e.DELETE("/api/stocks/{festival_stock_id}", fesStock.ID).
			Expect().
			Status(204)

		e.GET("/api/stocks/{festival_stock_id}", fesStock.ID).
			Expect().
			Status(404)
	})

	t.Run("Delete Festival Stock - Not Found", func(t *testing.T) {
		e.DELETE("/api/stocks/{festival_stock_id}", uuid.New()).
			Expect().
			Status(404)
	})

	t.Run("Delete Festival Stock - Invalid ID", func(t *testing.T) {
		e.DELETE("/api/stocks/{festival_stock_id}", "invalid-uuid").
			Expect().
			Status(404)
	})

	t.Run("Delete Festival Stock - Zero ID", func(t *testing.T) {
		e.DELETE("/api/stocks/{festival_stock_id}", uuid.Nil).
			Expect().
			Status(404)
	})
}