package v1

import (
	"log/slog"
	"time"

	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/router/utils/herror"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetVisitorCountsRequest struct {
	FestivalID string `param:"festival_id"`
}

func (h *Handler) GetVisitorCounts(c echo.Context) error {
	var req GetVisitorCountsRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request body")
	}

	fesID, err := uuid.Parse(req.FestivalID)
	if err != nil {
		return herror.BadRequest("Invalid festival ID format")
	}

	counts, err := h.r.ListVisitorCounts(c.Request().Context(), fesID)
	if err != nil {
		slog.Error("Failed to retrieve visitor counts", "error", err)
		switch err {
		case repository.ErrNotFound:
			return herror.NotFound("festival not found")
		default:
			return herror.InternalServerError("Failed to retrieve visitor counts")
		}
	}

	return c.JSON(200, map[string]any{
		"festival_id": fesID.String(),
		"counts":      counts,
	})
}

type AddVisitorCountRequest struct {
	FestivalID string `param:"festival_id"`
	Amount     int   `json:"amount" query:"amount"`
}

func (h *Handler) AddVisitorCount(c echo.Context) error {
	var req AddVisitorCountRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("Invalid request body")
	}

	fesID, err := uuid.Parse(req.FestivalID)
	if err != nil {
		return herror.BadRequest("Invalid festival ID format")
	}

	now := time.Now().UTC()

	err = h.r.AddVisitorCount(c.Request().Context(), fesID, now, req.Amount)
	if err != nil {
		slog.Error("Failed to add visitor count", "error", err)
		switch err {
		case repository.ErrNotFound:
			return herror.NotFound("festival not found")
		default:
			return herror.InternalServerError("Failed to add visitor count")
		}
	}

	return c.NoContent(204)
}
