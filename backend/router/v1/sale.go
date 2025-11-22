package v1

import (
	"log/slog"

	"github.com/Luke256/ducks/router/utils/herror"
	festivalstock "github.com/Luke256/ducks/service/festival_stock"
	"github.com/Luke256/ducks/service/sale"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateSaleRecordRequest struct {
	StockID  string `json:"stock_id"`
	Quantity int    `json:"quantity"`
}

func (r CreateSaleRecordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.StockID, validation.Required),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
	)
}

func (h *Handler) CreateSaleRecord(c echo.Context) error {
	var req CreateSaleRecordRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}
	stockID, err := uuid.Parse(req.StockID)
	if err != nil {
		return herror.NotFound("Stock not found")
	}

	record, err := h.saleManager.Create(stockID, req.Quantity)
	if err != nil {
		switch err {
		case festivalstock.ErrNotFound:
			return herror.NotFound("Stock not found")
		default:
			slog.Error("Failed to create sale record", "error", err)
			return herror.InternalServerError("Failed to create sale record")
		}
	}

	return c.JSON(201, record)
}

func (h *Handler) GetSaleRecord(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Sale record not found")
	}
	
	record, err := h.saleManager.Get(id)
	if err != nil {
		switch err {
		case sale.ErrNotFound:
			return herror.NotFound("Sale record not found")
		default:
			slog.Error("Failed to get sale record", "error", err)
			return herror.InternalServerError("Failed to get sale record")
		}
	}

	return c.JSON(200, record)
}

func (h *Handler) GetSaleRecordsByStockID(c echo.Context) error {
	stockID, err := uuid.Parse(c.Param("festival_stock_id"))
	if err != nil {
		return herror.NotFound("Sale records not found")
	}

	records, err := h.saleManager.GetByStockID(stockID)
	if err != nil {
		switch err {
		case festivalstock.ErrNotFound:
			return herror.NotFound("Sale records not found")
		default:
			slog.Error("Failed to get sale records by stock ID", "error", err)
			return herror.InternalServerError("Failed to get sale records")
		}
	}

	return c.JSON(200, map[string]any{"sales": records})
}

func (h *Handler) QuerySaleRecords(c echo.Context) error {
	festivalIDStr := c.QueryParam("festival_id")
	stockItemIDStr := c.QueryParam("stock_item_id")

	var festivalID, stockItemID uuid.UUID

	if festivalIDStr == "" {
		festivalID = uuid.Nil
	} else {
		id, err := uuid.Parse(festivalIDStr)
		if err != nil {
			return c.JSON(200, map[string]any{"sales": []sale.SaleRecord{}})
		}
		festivalID = id
	}

	if stockItemIDStr == "" {
		stockItemID = uuid.Nil
	} else {
		id, err := uuid.Parse(stockItemIDStr)
		if err != nil {
			return c.JSON(200, map[string]any{"sales": []sale.SaleRecord{}})
		}
		stockItemID = id
	}

	records, err := h.saleManager.Query(festivalID, stockItemID)
	if err != nil {
		slog.Error("Failed to query sale records", "error", err)
		return herror.InternalServerError("Failed to query sale records")
	}

	return c.JSON(200, map[string]any{"sales": records})
}

func (h *Handler) DeleteSaleRecord(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Sale record not found")
	}

	err = h.saleManager.Delete(id)
	if err != nil {
		switch err {
		case sale.ErrNotFound:
			return herror.NotFound("Sale record not found")
		default:
			slog.Error("Failed to delete sale record", "error", err)
			return herror.InternalServerError("Failed to delete sale record")
		}
	}

	return c.NoContent(204)
}
