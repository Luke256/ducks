package v1

import (
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	"github.com/Luke256/ducks/service/poster"
	"github.com/Luke256/ducks/utils/storage"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	r               repository.Repository
	festivalManager festival.Manager
	posterManager   poster.Manager
	storage         storage.Storage
}

func NewHandler(r repository.Repository, fm festival.Manager, pm poster.Manager, s storage.Storage) *Handler {
	return &Handler{
		r:               r,
		festivalManager: fm,
		posterManager:   pm,
		storage:         s,
	}
}

func (r *Handler) Setup(g *echo.Group) {
	festivals := g.Group("/festivals")
	posters := g.Group("/posters")
	images := g.Group("/images")

	// Images
	images.GET("/:id", r.GetImage)

	// Festivals
	festivals.POST("", r.CreateFestival)
	festivals.GET("", r.ListFestivals)
	festivals.GET("/:id", r.GetFestival)
	festivals.PUT("/:id", r.EditFestival)
	festivals.DELETE("/:id", r.DeleteFestival)

	// Posters
	posters.POST("", r.RegisterPoster)
	festivals.GET("/:festival_id/posters", r.ListPostersByFestival)
	posters.GET("/:id", r.GetPoster)
	posters.GET("/:festival_id/:poster_name", r.GetPosterByFestivalAndName)
	posters.PUT("/:id", r.EditPoster)
	posters.PATCH("/:id/status", r.UpdatePosterStatus)
	posters.DELETE("/:id", r.DeletePoster)
}
