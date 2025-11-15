package v1

import (
	"log/slog"

	"github.com/Luke256/ducks/service/poster"
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
	FestivalID  string `form:"festival_id" json:"festival_id"`
	PosterName  string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

func (r RegisterPosterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FestivalID, validation.Required),
		validation.Field(&r.PosterName, validation.Required, validation.Length(1, 64)),
		validation.Field(&r.Description, validation.Length(0, 1024)),
	)
}

type EditPosterRequest struct {
	ID          string `param:"id"`
	PosterName  string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

func (r EditPosterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.PosterName, validation.Required, validation.Length(1, 64)),
		validation.Field(&r.Description, validation.Length(0, 1024)),
	)
}

type UpdatePosterStatusRequest struct {
	ID     string `param:"id"`
	Status string `form:"status" json:"status"`
}

func (r UpdatePosterStatusRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Status, 
			validation.Required,
			validation.In(
				PosterStatusUncollected, 
				PosterStatusCollected, 
				PosterStatusLost,
			),
		),
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

	p, err := h.posterManager.Create(
		req.PosterName,
		fesID,
		req.Description,
		image,
	)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Festival not found")
		case poster.ErrAlreadyExists:
			return c.String(409, "Poster already exists")
		default:
			slog.Error("Failed to register poster", "error", err)
			return c.String(500, "Failed to register poster: "+ err.Error())
		}
	}

	return c.JSON(201, p)
}

func (h *Handler) ListPostersByFestival(c echo.Context) error {
	fesID, err := uuid.Parse(c.Param("festival_id"))
	if err != nil {
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
	posterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(404, "Poster not found")
	}

	p, err := h.posterManager.Get(posterID)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Poster not found")
		default:
			slog.Error("Failed to get poster", "error", err)
			return c.String(500, "Failed to get poster")
		}
	}

	return c.JSON(200, p)
}

func (h *Handler) GetPosterByFestivalAndName(c echo.Context) error {
	fesID, err := uuid.Parse(c.Param("festival_id"))
	if err != nil {
		return c.String(404, "Festival not found")
	}

	posterName := c.Param("poster_name")
	p, err := h.posterManager.GetByName(fesID, posterName)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Poster not found")
		default:
			slog.Error("Failed to get poster by festival and name", "error", err)
			return c.String(500, "Failed to get poster")
		}
	}

	return c.JSON(200, p)
}

func (h *Handler) EditPoster(c echo.Context) error {
	var req EditPosterRequest
	if err := c.Bind(&req); err != nil {
		return c.String(400, "Invalid request")
	}

	if err := req.Validate(); err != nil {
		return c.String(400, "Validation error: "+err.Error())
	}

	posterID, err := uuid.Parse(req.ID)
	if err != nil {
		return c.String(404, "Poster not found")
	}

	err = h.posterManager.Edit(posterID, req.PosterName, req.Description)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Poster not found")
		default:
			slog.Error("Failed to edit poster", "error", err)
			return c.String(500, "Failed to edit poster")
		}
	}

	return c.NoContent(204)
}

func (h *Handler) UpdatePosterStatus(c echo.Context) error {
	var req UpdatePosterStatusRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind update poster status request", "error", err)
		return c.String(400, "Invalid request")
	}

	if err := req.Validate(); err != nil {
		return c.String(400, "Validation error: "+err.Error())
	}

	posterID, err := uuid.Parse(req.ID)
	if err != nil {
		return c.String(404, "Poster not found")
	}

	err = h.posterManager.ChangeStatus(posterID, req.Status)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Poster not found")
		default:
			slog.Error("Failed to update poster status", "error", err)
			return c.String(500, "Failed to update poster status")
		}
	}

	return c.NoContent(204)
}

func (h *Handler) DeletePoster(c echo.Context) error {
	posterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(404, "Poster not found")
	}

	err = h.posterManager.Delete(posterID)
	if err != nil {
		switch err {
		case poster.ErrNotFound:
			return c.String(404, "Poster not found")
		default:
			slog.Error("Failed to delete poster", "error", err)
			return c.String(500, "Failed to delete poster")
		}
	}

	return c.NoContent(204)
}
