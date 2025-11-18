package v1

import (
	"github.com/Luke256/ducks/router/utils/herror"
	"github.com/Luke256/ducks/service/festival"
	festivalstock "github.com/Luke256/ducks/service/festival_stock"
	stockitem "github.com/Luke256/ducks/service/stock_item"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RegisterFestivalStockRequest struct {
	FestivalID  string `param:"festival_id"`
	StockItemID string `json:"item_id"`
	Price       int    `json:"price"`
}

func (r RegisterFestivalStockRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FestivalID, validation.Required),
		validation.Field(&r.StockItemID, validation.Required),
		validation.Field(&r.Price, validation.Required),
	)
}

type QueryFestivalStocksRequest struct {
	FestivalID string `param:"festival_id"`
	Category   string `query:"category"`
}

func (r QueryFestivalStocksRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FestivalID, validation.Required),
		validation.Field(&r.Category, validation.Length(1, 100)),
	)
}

type UpdateFestivalStockPriceRequest struct {
	ID    string `param:"id"`
	Price int    `json:"new_price"`
}

func (r UpdateFestivalStockPriceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Price, validation.Required, validation.Min(0)),
	)
}

func (h *Handler) RegisterFestivalStock(c echo.Context) error {
	var req RegisterFestivalStockRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}

	fesID, err := uuid.Parse(req.FestivalID)
	if err != nil {
		return herror.NotFound("Festival not found")
	}

	itemID, err := uuid.Parse(req.StockItemID)
	if err != nil {
		return herror.NotFound("Stock item not found")
	}

	festivalStock, err := h.festivalStockManager.Create(fesID, itemID, req.Price)
	if err != nil {
		switch err {
		case festival.ErrNotFound:
			return herror.NotFound("Festival not found")
		case stockitem.ErrNotFound:
			return herror.NotFound("Stock item not found")
		default:
			return herror.InternalServerError("Failed to create festival stock")
		}
	}

	return c.JSON(201, festivalStock)
}

func (h *Handler) GetFestivalStock(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Festival stock not found")
	}

	festivalStock, err := h.festivalStockManager.Get(id)
	if err != nil {
		switch err {
		case festivalstock.ErrNotFound:
			return herror.NotFound("Festival stock not found")
		default:
			return herror.InternalServerError("Failed to get festival stock")
		}
	}

	return c.JSON(200, festivalStock)
}

func (h *Handler) QueryFestivalStocks(c echo.Context) error {
	var req QueryFestivalStocksRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request parameters")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}

	fesID, err := uuid.Parse(req.FestivalID)
	if err != nil {
		return herror.NotFound("Festival not found")
	}

	festivalStocks, err := h.festivalStockManager.Query(fesID, req.Category)
	if err != nil {
		return herror.InternalServerError("Failed to query festival stocks")
	}

	return c.JSON(200, festivalStocks)
}

func (h *Handler) UpdateFestivalStockPrice(c echo.Context) error {
	var req UpdateFestivalStockPriceRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request body")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest("Validation failed: " + err.Error())
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		return herror.NotFound("Festival stock not found")
	}

	err = h.festivalStockManager.UpdatePrice(id, req.Price)
	if err != nil {
		switch err {
		case festivalstock.ErrNotFound:
			return herror.NotFound("Festival stock not found")
		default:
			return herror.InternalServerError("Failed to update festival stock price")
		}
	}

	return c.NoContent(204)
}

func (h *Handler) DeleteFestivalStock(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return herror.NotFound("Festival stock not found")
	}

	err = h.festivalStockManager.Delete(id)
	if err != nil {
		switch err {
		case festivalstock.ErrNotFound:
			return herror.NotFound("Festival stock not found")
		default:
			return herror.InternalServerError("Failed to delete festival stock")
		}
	}

	return c.NoContent(204)
}
