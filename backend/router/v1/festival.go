package v1

import (
	"github.com/Luke256/ducks/router/utils/herror"
	"github.com/Luke256/ducks/service/festival"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateFestivalRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CreateFestivalRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description),
	)
}

type EditFestivalRequest struct {
	ID          string `param:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r EditFestivalRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description),
	)
}

func (h *Handler) CreateFestival(c echo.Context) error {
	var req CreateFestivalRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("invalid request body")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest(err.Error())
	}

	festival, err := h.festivalManager.Create(req.Name, req.Description)
	if err != nil {
		return herror.InternalServerError("failed to create festival")
	}

	return c.JSON(201, festival)
}

func (h *Handler) ListFestivals(c echo.Context) error {
	festivals, err := h.festivalManager.List()
	if err != nil {
		return herror.InternalServerError("failed to list festivals")
	}
	return c.JSON(200, map[string]any{
		"festivals": festivals,
	})
}

func (h *Handler) GetFestival(c echo.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return herror.NotFound("festival not found")
	}

	fest, err := h.festivalManager.Get(id)
	if err != nil {
		switch err {
		case festival.ErrNotFound:
			return herror.NotFound("festival not found")
		default:
			return herror.InternalServerError("failed to get festival")
		}
	}
	return c.JSON(200, fest)
}

func (h *Handler) EditFestival(c echo.Context) error {
	var req EditFestivalRequest
	if err := c.Bind(&req); err != nil {
		return herror.BadRequest("invalid request body")
	}
	if err := req.Validate(); err != nil {
		return herror.BadRequest(err.Error())
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return herror.NotFound("festival not found")
	}

	err = h.festivalManager.Edit(id, req.Name, req.Description)
	if err != nil {
		switch err {
		case festival.ErrNotFound:
			return herror.NotFound("festival not found")
		default:
			return herror.InternalServerError("failed to edit festival")
		}
	}

	fest, err := h.festivalManager.Get(id)
	if err != nil {
		return herror.InternalServerError("failed to get festival after edit")
	}

	return c.JSON(200, fest)
}

func (h *Handler) DeleteFestival(c echo.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return herror.NotFound("festival not found")
	}

	err = h.festivalManager.Delete(id)
	if err != nil {
		switch err {
		case festival.ErrNotFound:
			return herror.NotFound("festival not found")
		default:
			return herror.InternalServerError("failed to delete festival")
		}
	}

	return c.NoContent(204)
}
