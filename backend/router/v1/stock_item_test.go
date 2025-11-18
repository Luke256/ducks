package v1

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)


func TestRegisterStockItem(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	t.Run("RegisterStockItem", func(t *testing.T) {
		req := e.POST("/api/stocks").
			WithMultipart().
			WithForm(map[string]any{
				"name":  "Sample Stock Item",
				"description": "This is a sample stock item.",
				"category": "Sample Category",
			}).
			WithFile("image", "sample_image.png", strings.NewReader("")).
			Expect().
			Status(201).
			JSON().
			Object()

		req.Value("name").IsEqual("Sample Stock Item")
		req.Value("description").IsEqual("This is a sample stock item.")
		req.Value("category").IsEqual("Sample Category")
		req.Value("image_url").NotNull()
	})

	t.Run("Empty Name", func(t *testing.T) {
		e.POST("/api/stocks").
			WithMultipart().
			WithForm(map[string]any{
				"name":  "",
				"description": "This is a sample stock item.",
				"category": "Sample Category",
			}).
			WithFile("image", "sample_image.png", strings.NewReader("")).
			Expect().
			Status(400)
	})

	t.Run("Missing Image", func(t *testing.T) {
		e.POST("/api/stocks").
			WithMultipart().
			WithForm(map[string]any{
				"name":  "Sample Stock Item",
				"description": "This is a sample stock item.",
				"category": "Sample Category",
			}).
			Expect().
			Status(400)
	})

	t.Run("Empty Category", func(t *testing.T) {
		e.POST("/api/stocks").
			WithMultipart().
			WithForm(map[string]any{
				"name":  "Sample Stock Item",
				"description": "This is a sample stock item.",
				"category": "",
			}).
			WithFile("image", "sample_image.png", strings.NewReader("")).
			Expect().
			Status(400)
	})
}

func TestGetStockItem(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	item := env.mustCreateStockItem(t, "Test Item", "This is a test item.", "Test Category")

	t.Run("GetStockItem", func(t *testing.T) {
		resp := e.GET("/api/stocks/{id}", item.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("id").IsEqual(item.ID.String())
		resp.Value("name").IsEqual(item.Name)
		resp.Value("description").IsEqual(item.Description)
		resp.Value("category").IsEqual(item.Category)
		resp.Value("image_url").NotNull()
	})

	t.Run("GetStockItem Not Found", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.GET("/api/stocks/{id}", id.String()).
			Expect().
			Status(404)
	})

	t.Run("GetStockItem Zero ID", func(t *testing.T) {
		e.GET("/api/stocks/{id}", "00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})

	t.Run("GetStockItem Invalid ID", func(t *testing.T) {
		e.GET("/api/stocks/{id}", "invalid-uuid").
			Expect().
			Status(404)
	})
}

func TestQueryStockItems(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	item1 := env.mustCreateStockItem(t, "Item 1", "Description 1", "Category A")
	item2 := env.mustCreateStockItem(t, "Item 2", "Description 2", "Category A")
	item3 := env.mustCreateStockItem(t, "Item 3", "Description 3", "Category B")

	t.Run("QueryStockItems", func(t *testing.T) {
		resp := e.GET("/api/stocks").
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("items").Array().ContainsOnly(item1, item2, item3)
	})

	t.Run("QueryStockItems with Category Filter", func(t *testing.T) {
		resp := e.GET("/api/stocks").
			WithQuery("category", "Category A").
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("items").Array().ContainsOnly(item1, item2)
	})
}

func TestEditStockItem(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	item := env.mustCreateStockItem(t, "Original Name", "Original Description", "Original Category")

	t.Run("EditStockItem", func(t *testing.T) {
		e.PUT("/api/stocks/{id}", item.ID.String()).
			WithJSON(map[string]any{
				"name":        "Updated Name",
				"description": "Updated Description",
				"category":    "Updated Category",
			}).
			Expect().
			Status(204)
		
		res := e.GET("/api/stocks/{id}", item.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("id").IsEqual(item.ID.String())
		res.Value("name").IsEqual("Updated Name")
		res.Value("description").IsEqual("Updated Description")
		res.Value("category").IsEqual("Updated Category")
	})

	t.Run("EditStockItem Not Found", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.PUT("/api/stocks/{id}", id.String()).
			WithJSON(map[string]any{
				"name":        "Updated Name",
				"description": "Updated Description",
				"category":    "Updated Category",
			}).
			Expect().
			Status(404)
	})

	t.Run("EditStockItem Zero ID", func(t *testing.T) {
		e.PUT("/api/stocks/{id}", "00000000-0000-0000-0000-000000000000").
			WithJSON(map[string]any{
				"name":        "Updated Name",
				"description": "Updated Description",
				"category":    "Updated Category",
			}).
			Expect().
			Status(404)
	})

	t.Run("EditStockItem Invalid ID", func(t *testing.T) {
		e.PUT("/api/stocks/{id}", "invalid-uuid").
			WithJSON(map[string]any{
				"name":        "Updated Name",
				"description": "Updated Description",
				"category":    "Updated Category",
			}).
			Expect().
			Status(404)
	})
}

func TestUpdateStockItemImage(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	item := env.mustCreateStockItem(t, "Item with Image", "Description", "Category")

	t.Run("UpdateStockItemImage", func(t *testing.T) {
		e.PUT("/api/stocks/{id}/image", item.ID.String()).
			WithMultipart().
			WithFile("image", "new_image.png", strings.NewReader("")).
			Expect().
			Status(204)
	})

	t.Run("UpdateStockItemImage Not Found", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}

		e.PUT("/api/stocks/{id}/image", id.String()).
			WithMultipart().
			WithFile("image", "new_image.png", strings.NewReader("")).
			Expect().
			Status(404)
	})

	t.Run("UpdateStockItemImage Zero ID", func(t *testing.T) {
		e.PUT("/api/stocks/{id}/image", "00000000-0000-0000-0000-000000000000").
			WithMultipart().
			WithFile("image", "new_image.png", strings.NewReader("")).
			Expect().
			Status(404)
	})

	t.Run("UpdateStockItemImage Invalid ID", func(t *testing.T) {
		e.PUT("/api/stocks/{id}/image", "invalid-uuid").
			WithMultipart().
			WithFile("image", "new_image.png", strings.NewReader("")).
			Expect().
			Status(404)
	})
}

func TestDeleteStockItem(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	item := env.mustCreateStockItem(t, "Item to Delete", "Description", "Category")

	t.Run("DeleteStockItem", func(t *testing.T) {
		e.DELETE("/api/stocks/{id}", item.ID.String()).
			Expect().
			Status(204)

		e.GET("/api/stocks/{id}", item.ID.String()).
			Expect().
			Status(404)
	})

	t.Run("DeleteStockItem Not Found", func(t *testing.T) {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatalf("failed to generate uuid: %v", err)
		}
		e.DELETE("/api/stocks/{id}", id.String()).
			Expect().
			Status(404)
	})

	t.Run("DeleteStockItem Zero ID", func(t *testing.T) {
		e.DELETE("/api/stocks/{id}", "00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})

	t.Run("DeleteStockItem Invalid ID", func(t *testing.T) {
		e.DELETE("/api/stocks/{id}", "invalid-uuid").
			Expect().
			Status(404)
	})
}