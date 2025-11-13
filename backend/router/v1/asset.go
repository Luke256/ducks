package v1

import (
	"log/slog"

	"github.com/Luke256/ducks/router/utils/herror"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetImage(c echo.Context) error {
	imageID := c.Param("id")

	file, err := h.storage.DownloadFile(imageID)
	if err != nil {
		slog.Error("failed to download image", "error", err, "image_id", imageID)
		return herror.NotFound()
	}

	c.Response().Header().Set(echo.HeaderCacheControl, "max-age=31536000, immutable")

	return c.Stream(200, "image/webp", file)
}