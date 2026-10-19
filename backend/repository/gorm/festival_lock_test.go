package gorm

import (
	"testing"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLockedFestivalRejectsMutations(t *testing.T) {
	t.Parallel()
	r := setup(t, common)
	festival := mustCreateFestival(t, r, "Locked event", "description")
	item := mustCreateStockItem(t, r, "Item", "description", "category", uuid.NewString())
	stock := mustCreateFestivalStock(t, r, festival.ID, item.ID, 100, "description", true)
	poster := mustCreatePoster(t, r, festival.ID, "Poster", "description", uuid.NewString())
	sale := mustCreateSaleRecord(t, r, stock.ID, 2)
	now := time.Now().UTC()
	require.NoError(t, r.AddVisitorCount(t.Context(), festival.ID, now, 5))

	otherFestival := mustCreateFestival(t, r, "Open event", "description")
	otherStock := mustCreateFestivalStock(t, r, otherFestival.ID, item.ID, 100, "description", true)
	require.NoError(t, r.db.Model(&festival).Update("locked", true).Error)

	for _, tc := range []struct {
		name string
		run  func() error
	}{
		{"edit event", func() error {
			return r.UpdateFestival(festival.ID, "changed", "changed", festival.Date.AddDate(0, 0, 1))
		}},
		{"delete event", func() error { return r.DeleteFestival(festival.ID) }},
		{"create poster", func() error {
			_, err := r.RegisterPoster(festival.ID, "New poster", "changed", []string{uuid.NewString()})
			return err
		}},
		{"edit poster", func() error { return r.UpdatePoster(poster.ID, "changed", "changed") }},
		{"change poster status", func() error { return r.UpdatePosterStatus(poster.ID, "collected") }},
		{"change poster images", func() error {
			_, err := r.UpdatePosterImages(poster.ID, []string{uuid.NewString()}, []string{poster.Images[0].ID})
			return err
		}},
		{"delete poster", func() error { return r.DeletePoster(poster.ID) }},
		{"create stock", func() error {
			_, err := r.RegisterFestivalStock(festival.ID, item.ID, 200, "changed", true)
			return err
		}},
		{"edit stock", func() error { return r.UpdateFestivalStock(stock.ID, "changed", false) }},
		{"delete stock", func() error { return r.DeleteFestivalStock(stock.ID) }},
		{"create mixed sales", func() error {
			_, err := r.CreateSaleRecords(
				repository.SaleData{FestivalStockID: otherStock.ID, Quantity: 1},
				repository.SaleData{FestivalStockID: stock.ID, Quantity: 1},
			)
			return err
		}},
		{"delete sale", func() error { return r.DeleteSaleRecord(sale.ID) }},
		{"increment visitors", func() error { return r.AddVisitorCount(t.Context(), festival.ID, now, 1) }},
		{"decrement visitors", func() error { return r.AddVisitorCount(t.Context(), festival.ID, now, -1) }},
		{"edit shared item", func() error {
			_, err := r.UpdateStockItem(item.ID, "changed", "changed", "changed", "changed")
			return err
		}},
		{"delete shared item", func() error { return r.DeleteStockItem(item.ID) }},
	} {
		t.Run(tc.name, func(t *testing.T) { require.ErrorIs(t, tc.run(), repository.ErrEventLocked) })
	}

	gotFestival, err := r.GetFestivalByID(festival.ID)
	require.NoError(t, err)
	festival.Locked = true
	require.Equal(t, festival, gotFestival)
	gotPoster, err := r.GetPosterByID(poster.ID)
	require.NoError(t, err)
	require.Equal(t, poster.PosterName, gotPoster.PosterName)
	require.Equal(t, poster.Description, gotPoster.Description)
	require.Equal(t, poster.Status, gotPoster.Status)
	require.Equal(t, poster.Images, gotPoster.Images)
	gotStock, err := r.GetFestivalStockByID(stock.ID)
	require.NoError(t, err)
	require.Equal(t, stock.Description, gotStock.Description)
	require.Equal(t, stock.ForSale, gotStock.ForSale)
	gotItem, err := r.GetStockItemByID(item.ID)
	require.NoError(t, err)
	require.Equal(t, item, gotItem)
	sales, err := r.GetSaleRecordsByFestivalStockID(stock.ID)
	require.NoError(t, err)
	require.Len(t, sales, 1)
	require.Equal(t, sale.ID, sales[0].ID)
	otherSales, err := r.GetSaleRecordsByFestivalStockID(otherStock.ID)
	require.NoError(t, err)
	require.Empty(t, otherSales)
	counts, err := r.ListVisitorCounts(t.Context(), festival.ID)
	require.NoError(t, err)
	require.Len(t, counts, 1)
	require.EqualValues(t, 5, counts[0].Count)

	require.NoError(t, r.db.Model(&festival).Update("locked", false).Error)
	require.NoError(t, r.UpdateFestival(festival.ID, "Unlocked event", "updated", festival.Date))
	require.NoError(t, r.AddVisitorCount(t.Context(), festival.ID, now, 1))
}

func TestMutationWaitsForConcurrentFestivalLock(t *testing.T) {
	t.Parallel()
	r := setup(t, common)
	festival := mustCreateFestival(t, r, "Event", "description")
	tx := r.db.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { tx.Rollback() })
	require.NoError(t, tx.Model(&model.Festival{}).Where("id = ?", festival.ID).Update("locked", true).Error)

	done := make(chan error, 1)
	go func() { done <- r.UpdateFestival(festival.ID, "changed", "changed", festival.Date.AddDate(0, 0, 1)) }()
	select {
	case err := <-done:
		t.Fatalf("mutation completed before lock committed: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	require.NoError(t, tx.Commit().Error)
	select {
	case err := <-done:
		require.ErrorIs(t, err, repository.ErrEventLocked)
	case <-time.After(5 * time.Second):
		t.Fatal("mutation did not finish after lock committed")
	}
	got, err := r.GetFestivalByID(festival.ID)
	require.NoError(t, err)
	require.Equal(t, festival.Name, got.Name)
	require.True(t, festival.Date.Equal(got.Date))
}
