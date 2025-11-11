package v1

import (
	"github.com/Luke256/ducks/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	r repository.Repository
}

func NewHandler(r repository.Repository) *Handler {
	return &Handler{
		r: r,
	}
}

func (r *Handler) Setup(g *echo.Group) {
	
}