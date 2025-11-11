package v1

import (
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	"github.com/Luke256/ducks/service/poster"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	r repository.Repository
	festivalManager festival.Manager
	posterManager   poster.Manager
}

func NewHandler(r repository.Repository, fm festival.Manager, pm poster.Manager) *Handler {
	return &Handler{
		r: r,
		festivalManager: fm,
		posterManager: pm,
	}
}

func (r *Handler) Setup(g *echo.Group) {
	
}