package repository

import (
	"time"
	"context"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
)

type VisitorCounterRepository interface {
	// ListVisitorCounts は指定したイベントの来客数を取得します
	ListVisitorCounts(ctx context.Context, festivalID uuid.UUID) ([]model.VisitorCount, error)
	
	// AddVisitorCount は指定したイベントの来客数を加算します
	AddVisitorCount(ctx context.Context, festivalID uuid.UUID, timestamp time.Time, amount int) error
}
