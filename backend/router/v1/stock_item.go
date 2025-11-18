package v1

import (
	"log/slog"

	"github.com/Luke256/ducks/router/utils/herror"
	stockitem "github.com/Luke256/ducks/service/stock_item"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RegisterStockItemRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Category    string `form:"category"`
}

func (r RegisterStockItemRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Category, validation.Required, validation.Length(1, 100)),
	)
}

func (h *Handler) RegisterStockItem(c echo.Context) error {
	var req RegisterStockItemRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}

	image, err := c.FormFile("image")
	if err != nil {
		return herror.BadRequest("Image is required")
	}

	item, err := h.stockItemManager.Create(req.Name, req.Description, req.Category, image)
	if err != nil {
		slog.Error("failed to register stock item:", slog.String("error", err.Error()))
		return c.String(500, "Failed to register stock item")
	}

	return c.JSON(201, item)
}

func (h *Handler) GetStockItems(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Stock item not found")
	}

	item, err := h.stockItemManager.Get(id)
	if err != nil {
		switch err {
		case stockitem.ErrNotFound:
			return herror.NotFound("Stock item not found")
		default:
			slog.Error("failed to get stock item:", slog.String("error", err.Error()))
			return c.String(500, "Failed to get stock item")
		}
	}

	return c.JSON(200, item)
}

func (h *Handler) QueryStockItems(c echo.Context) error {
	category := c.QueryParam("category")

	items, err := h.stockItemManager.Query(category)
	if err != nil {
		slog.Error("failed to query stock items:", slog.String("error", err.Error()))
		return c.String(500, "Failed to query stock items")
	}

	return c.JSON(200, map[string]any{
		"items": items,
	})
}

func (h *Handler) EditStockItem(c echo.Context) error {
	var req RegisterStockItemRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Stock item not found")
	}

	err = h.stockItemManager.Edit(id, req.Name, req.Description, req.Category)
	if err != nil {
		switch err {
		case stockitem.ErrNotFound:
			return herror.NotFound("Stock item not found")
		default:
			slog.Error("failed to edit stock item:", slog.String("error", err.Error()))
			return c.String(500, "Failed to edit stock item")
		}
	}

	return c.NoContent(204)
}

func (h *Handler) UpdateStockItemImage(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Stock item not found")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return herror.BadRequest("Image is required")
	}

	err = h.stockItemManager.UpdateImage(id, image)
	if err != nil {
		switch err {
		case stockitem.ErrNotFound:
			return herror.NotFound("Stock item not found")
		default:
			slog.Error("failed to update stock item image:", slog.String("error", err.Error()))
			return c.String(500, "Failed to update stock item image")
		}
	}

	return c.NoContent(204)
}

func (h *Handler) DeleteStockItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Stock item not found")
	}

	err = h.stockItemManager.Delete(id)
	if err != nil {
		switch err {
		case stockitem.ErrNotFound:
			return herror.NotFound("Stock item not found")
		default:
			slog.Error("failed to delete stock item:", slog.String("error", err.Error()))
			return c.String(500, "Failed to delete stock item")
		}
	}
	return c.NoContent(204)
}