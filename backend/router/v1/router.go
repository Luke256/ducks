package v1

import (
	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/festival"
	festivalstock "github.com/Luke256/ducks/service/festival_stock"
	"github.com/Luke256/ducks/service/poster"
	"github.com/Luke256/ducks/service/sale"
	stockitem "github.com/Luke256/ducks/service/stock_item"
	"github.com/Luke256/ducks/utils/storage"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	r                    repository.Repository
	festivalManager      festival.Manager
	posterManager        poster.Manager
	stockItemManager     stockitem.Manager
	festivalStockManager festivalstock.Manager
	saleManager          sale.Manager
	storage              storage.Storage
}

func NewHandler(r repository.Repository, fm festival.Manager, pm poster.Manager, sim stockitem.Manager, fsm festivalstock.Manager, sm sale.Manager, s storage.Storage) *Handler {
	return &Handler{
		r:                    r,
		festivalManager:      fm,
		posterManager:        pm,
		stockItemManager:     sim,
		festivalStockManager: fsm,
		saleManager:          sm,
		storage:              s,
	}
}

func (r *Handler) Setup(g *echo.Group) {
	festivals := g.Group("/festivals")
	posters := g.Group("/posters")
	images := g.Group("/images")
	stockItems := g.Group("/stocks")
	festivalStocks := g.Group("/festival_stocks")
	sales := g.Group("/sales")

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

	// Stock Items
	stockItems.POST("", r.RegisterStockItem)
	stockItems.GET("", r.QueryStockItems)
	stockItems.GET("/:id", r.GetStockItems)
	stockItems.PUT("/:id", r.EditStockItem)
	stockItems.PUT("/:id/image", r.UpdateStockItemImage)
	stockItems.DELETE("/:id", r.DeleteStockItem)

	// Festival Stocks
	festivals.POST("/:festival_id/stocks", r.RegisterFestivalStock)
	festivals.GET("/:festival_id/stocks", r.QueryFestivalStocks)
	festivalStocks.GET("/:id", r.GetFestivalStock)
	festivalStocks.PUT("/:id/price", r.UpdateFestivalStockPrice)
	festivalStocks.DELETE("/:id", r.DeleteFestivalStock)

	// Sales
	sales.POST("", r.CreateSaleRecord)
	sales.GET("/:id", r.GetSaleRecord)
	festivalStocks.GET("/:festival_stock_id/sales", r.GetSaleRecordsByStockID)
	sales.GET("", r.QuerySaleRecords)
	sales.DELETE("/:id", r.DeleteSaleRecord)
}
