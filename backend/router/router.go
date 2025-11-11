package router

import (
	"github.com/Luke256/ducks/router/v1"
	"github.com/Luke256/ducks/repository"

	"github.com/labstack/echo/v4"
)

type Router struct {
	e *echo.Echo
	r repository.Repository
	v1 *v1.Handler
}

func NewRouter(e *echo.Echo, v1Handler *v1.Handler, repo repository.Repository) *Router {
	return &Router{
		e: e,
		r: repo,
		v1: v1Handler,
	}
}

func (r *Router) Setup() {
	api := r.e.Group("/api")

	v1 := api.Group("/v1")
	r.v1.Setup(v1)
}

func (r *Router) Start() error {
	return r.e.Start(":8080")
}