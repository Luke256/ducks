package migration

import (
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestV4MigratesExistingStocksForSale(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.FestivalStock{}))
	require.NoError(t, db.Migrator().DropColumn(&model.FestivalStock{}, "ForSale"))
	festival := model.Festival{ID: uuid.New(), Name: "test"}
	item := model.StockItem{ID: uuid.New(), Name: "test"}
	require.NoError(t, db.Create(&festival).Error)
	require.NoError(t, db.Create(&item).Error)
	id := uuid.New()
	require.NoError(t, db.Exec(
		"INSERT INTO festival_stocks (id, festival_id, stock_item_id, price, description) VALUES (?, ?, ?, ?, ?)",
		id, festival.ID, item.ID, 100, "Existing stock",
	).Error)

	require.NoError(t, v4().Migrate(db))
	var stock model.FestivalStock
	require.NoError(t, db.First(&stock, "id = ?", id).Error)
	require.True(t, stock.ForSale)
	require.Equal(t, 100, stock.Price)
	require.Equal(t, "Existing stock", stock.Description)

	// 再実行しても、移行後に非表示にした商品を販売中へ戻さない。
	require.NoError(t, db.Model(&stock).Update("for_sale", false).Error)
	require.NoError(t, v4().Migrate(db))
	require.NoError(t, db.First(&stock, "id = ?", id).Error)
	require.False(t, stock.ForSale)
}
