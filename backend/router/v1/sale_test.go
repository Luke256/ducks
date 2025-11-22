package v1

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateSaleRecord(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "Description")
	stock_item := env.mustCreateStockItem(t, "Test Stock Item", "Category", "")
	stock := env.mustCreateFestivalStock(t, fes.ID, stock_item.ID, 100)

	t.Run("Create Sale Record", func(t *testing.T) {
		res := e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": stock.ID.String(),
				"quantity": 2,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		res.Value("stock_id").IsEqual(stock.ID.String())
		res.Value("quantity").IsEqual(2)
	})

	t.Run("Invalid Stock ID", func(t *testing.T) {
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": "invalid-uuid",
				"quantity": 2,
			}).
			Expect().
			Status(404)
	})

	t.Run("Unexisting Stock ID", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": id.String(),
				"quantity": 2,
			}).
			Expect().
			Status(404)
	})

	t.Run("Missing Stock ID", func(t *testing.T) {
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"quantity": 2,
			}).
			Expect().
			Status(400)
	})

	t.Run("Negative Quantity", func(t *testing.T) {
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": stock.ID.String(),
				"quantity": -5,
			}).
			Expect().
			Status(400)
	})

	t.Run("Zero Quantity", func(t *testing.T) {
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": stock.ID.String(),
				"quantity": 0,
			}).
			Expect().
			Status(400)
	})

	t.Run("Missing Quantity", func(t *testing.T) {
		e.POST("/api/sales").
			WithJSON(map[string]any{
				"stock_id": stock.ID.String(),
			}).
			Expect().
			Status(400)
	})
}

func TestGetSaleRecord(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "Description")
	stock_item := env.mustCreateStockItem(t, "Test Stock Item", "Category", "")
	stock := env.mustCreateFestivalStock(t, fes.ID, stock_item.ID, 100)
	record := env.mustCreateSaleRecord(t, stock.ID, 3)

	t.Run("Get Existing Sale Record", func(t *testing.T) {
		res := e.GET("/api/sales/{id}", record.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("id").IsEqual(record.ID.String())
		res.Value("stock_id").IsEqual(stock.ID)
		res.Value("quantity").IsEqual(3)
	})

	t.Run("Get Non-Existing Sale Record", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.GET("/api/sales/{id}", id.String()).
			Expect().
			Status(404)
	})

	t.Run("Invalid Sale Record ID", func(t *testing.T) {
		e.GET("/api/sales/{id}", "invalid-uuid").
			Expect().
			Status(404)
	})
}

func TestGetSaleRecordsByStockID(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "Description")
	stock_item := env.mustCreateStockItem(t, "Test Stock Item", "Category", "")
	stock1 := env.mustCreateFestivalStock(t, fes.ID, stock_item.ID, 100)
	stock2 := env.mustCreateFestivalStock(t, fes.ID, stock_item.ID, 150)

	sale1 := env.mustCreateSaleRecord(t, stock1.ID, 2)
	sale2 := env.mustCreateSaleRecord(t, stock1.ID, 5)
	env.mustCreateSaleRecord(t, stock2.ID, 3)

	t.Run("Get Sale Records by Stock ID", func(t *testing.T) {
		res := e.GET("/api/festival_stocks/{stock_id}/sales", stock1.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("sales").Array().ContainsOnly(sale1, sale2)
	})

	t.Run("Get Sale Records by Non-Existing Stock ID", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.GET("/api/festival_stocks/{stock_id}/sales", id).
			Expect().
			Status(404)
	})

	t.Run("Invalid Stock ID", func(t *testing.T) {
		e.GET("/api/festival_stocks/{stock_id}/sales", "invalid-uuid").
			Expect().
			Status(404)
	})
}

func TestQuerySaleRecords(t *testing.T) {
	env := setup(t, s3)
	e := env.R(t)

	fes1 := env.mustCreateFestival(t, "Festival 1", "Description 1")
	fes2 := env.mustCreateFestival(t, "Festival 2", "Description 2")
	stock_item1 := env.mustCreateStockItem(t, "Stock Item 1", "Category", "")
	stock_item2 := env.mustCreateStockItem(t, "Stock Item 2", "Category", "")

	stock1 := env.mustCreateFestivalStock(t, fes1.ID, stock_item1.ID, 100)
	stock2 := env.mustCreateFestivalStock(t, fes1.ID, stock_item2.ID, 150)
	stock3 := env.mustCreateFestivalStock(t, fes2.ID, stock_item1.ID, 200)

	sale1 := env.mustCreateSaleRecord(t, stock1.ID, 2)
	sale2 := env.mustCreateSaleRecord(t, stock1.ID, 3)
	sale3 := env.mustCreateSaleRecord(t, stock2.ID, 5)
	sale4 := env.mustCreateSaleRecord(t, stock3.ID, 4)

	t.Run("Query All Sale Records", func(t *testing.T) {
		res := e.GET("/api/sales").
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("sales").Array().ContainsOnly(sale1, sale2, sale3, sale4)
	})

	t.Run("Query Sale Records by Festival ID", func(t *testing.T) {
		res := e.GET("/api/sales").
			WithQuery("festival_id", fes1.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("sales").Array().ContainsOnly(sale1, sale2, sale3)
	})

	t.Run("Query Sale Records by Stock Item ID", func(t *testing.T) {
		res := e.GET("/api/sales").
			WithQuery("stock_item_id", stock_item1.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("sales").Array().ContainsOnly(sale1, sale2, sale4)
	})

	t.Run("Query Sale Records by Festival ID and Stock Item ID", func(t *testing.T) {
		res := e.GET("/api/sales").
			WithQuery("festival_id", fes1.ID.String()).
			WithQuery("stock_item_id", stock_item1.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("sales").Array().ContainsOnly(sale1, sale2)
	})

	t.Run("Query Sale Records with Invalid Festival ID", func(t *testing.T) {
		res := e.GET("/api/sales").
			WithQuery("festival_id", "invalid-uuid").
			Expect().
			Status(200).
			JSON().Object()

		res.Value("sales").Array().IsEmpty()
	})

	t.Run("Query Sale Records with Invalid Stock Item ID", func(t *testing.T) {
		res := e.GET("/api/sales").
			WithQuery("stock_item_id", "invalid-uuid").
			Expect().
			Status(200).
			JSON().Object()

		res.Value("sales").Array().IsEmpty()
	})
}

func TestDeleteSaleRecord(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Test Festival", "Description")
	stock_item := env.mustCreateStockItem(t, "Test Stock Item", "Category", "")
	stock := env.mustCreateFestivalStock(t, fes.ID, stock_item.ID, 100)

	sale := env.mustCreateSaleRecord(t, stock.ID, 4)

	t.Run("Delete Sale Record", func (t *testing.T) {
		e.DELETE("/api/sales/{id}", sale.ID).
			Expect().
			Status(204)

		e.GET("/api/sales/{id}", sale.ID).
			Expect().
			Status(404)
	})

	t.Run("Delete Non-Existing Sale Record", func (t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.DELETE("/api/sales/{id}", id).
			Expect().
			Status(404)
	})

	t.Run("Delete Sale Record with Invalid ID", func (t *testing.T) {
		e.DELETE("/api/sales/{id}", "invalid-uuid").
			Expect().
			Status(404)
	})

	t.Run("Delete Sale Record with Zero UUID", func (t *testing.T) {
		e.DELETE("/api/sales/{id}", uuid.Nil).
			Expect().
			Status(404)
	})
}