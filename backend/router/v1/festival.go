package v1

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateFestival(c echo.Context) error {
	return c.String(200, "CreateFestival")
}

func (h *Handler) ListFestivals(c echo.Context) error {
	return c.String(200, "ListFestivals")
}

func (h *Handler) GetFestival(c echo.Context) error {
	return c.String(200, "GetFestival")
}

func (h *Handler) EditFestival(c echo.Context) error {
	return c.String(200, "EditFestival")
}

func (h *Handler) DeleteFestival(c echo.Context) error {
	return c.String(200, "DeleteFestival")
}