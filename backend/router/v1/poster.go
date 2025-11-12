package v1

import (
	"log/slog"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	PosterStatusUncollected = "uncollected"
	PosterStatusCollected   = "collected"
	PosterStatusLost        = "lost"
)

type RegisterPosterRequest struct {
	FestivalID  string                `form:"festival_id"`
	PosterName  string                `form:"name"`
	Description string                `form:"description"`
}

func (r RegisterPosterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FestivalID, validation.Required),
		validation.Field(&r.PosterName, validation.Required, validation.Length(1, 64)),
		validation.Field(&r.Description, validation.Length(0, 1024)),
	)
}

func (h *Handler) RegisterPoster(c echo.Context) error {
	var req RegisterPosterRequest
	if err := c.Bind(&req); err != nil {
		return c.String(400, "Invalid request")
	}

	if err := req.Validate(); err != nil {
		return c.String(400, "Validation error: "+err.Error())
	}

	fesID, err := uuid.Parse(req.FestivalID)
	if err != nil {
		return c.String(404, "Festival not found")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.String(400, "Invalid image file: "+err.Error())
	}

	poster, err := h.posterManager.Create(
		req.PosterName,
		fesID,
		req.Description,
		image,
	)
	if err != nil {
		slog.Error("Failed to create poster", "error", err)
		return c.String(500, "Failed to create poster")
	}

	return c.JSON(201, poster)
}

func (h *Handler) ListPostersByFestival(c echo.Context) error {
	fesID, err := uuid.Parse(c.Param("festival_id"))
	if err != nil {
		slog.Error("Invalid festival ID", "error", err)
		return c.String(404, "Festival not found")
	}

	posters, err := h.posterManager.GetByFestival(fesID)
	if err != nil {
		slog.Error("Failed to list posters by festival", "error", err)
		return c.String(500, "Failed to list posters")
	}

	return c.JSON(200, map[string]any{
		"posters": posters,
	})
}

func (h *Handler) GetPoster(c echo.Context) error {
	return c.String(200, "GetPoster")
}

func (h *Handler) GetPosterByFestivalAndName(c echo.Context) error {
	return c.String(200, "GetPosterByFestivalAndName")
}

func (h *Handler) EditPoster(c echo.Context) error {
	return c.String(200, "EditPoster")
}

func (h *Handler) UpdatePosterStatus(c echo.Context) error {
	return c.String(200, "UpdatePosterStatus")
}

func (h *Handler) DeletePoster(c echo.Context) error {
	return c.String(200, "DeletePoster")
}
