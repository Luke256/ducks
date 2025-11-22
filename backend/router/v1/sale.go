package v1

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateSaleRecord(c echo.Context) error {
	return c.String(200, "CreateSaleRecord")
}

func (h *Handler) GetSaleRecord(c echo.Context) error {
	return c.String(200, "GetSaleRecord")
}

func (h *Handler) GetSaleRecordsByStockID(c echo.Context) error {
	return c.String(200, "GetSaleRecordsByStockID")
}

func (h *Handler) QuerySaleRecords(c echo.Context) error {
	return c.String(200, "QuerySaleRecords")
}

func (h *Handler) DeleteSaleRecord(c echo.Context) error {
	return c.String(200, "DeleteSaleRecord")
}