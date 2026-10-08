package gorm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVisitorCountUTC(t *testing.T) {
	repo := setup(t, common)
	festival := mustCreateFestival(t, repo, "UTC Visitor Festival", "UTC test")
	timestamp := time.Date(2026, 10, 8, 0, 9, 59, 0, time.FixedZone("JST", 9*60*60))
	require.NoError(t, repo.AddVisitorCount(t.Context(), festival.ID, timestamp, 5))
	// The same instant with a different offset must use the same UTC bucket.
	require.NoError(t, repo.AddVisitorCount(t.Context(), festival.ID, timestamp.UTC(), -3))
	counts, err := repo.ListVisitorCounts(t.Context(), festival.ID)
	require.NoError(t, err)
	require.Len(t, counts, 1)
	require.Equal(t, time.Date(2026, 10, 7, 15, 0, 0, 0, time.UTC), counts[0].BucketStart)
	require.Equal(t, uint(2), counts[0].Count)
	var stored string
	err = repo.db.Raw("SELECT DATE_FORMAT(bucket_start, '%Y-%m-%d %H:%i:%s') FROM visitor_counts WHERE festival_id = ?", festival.ID).Scan(&stored).Error
	require.NoError(t, err)
	require.Equal(t, "2026-10-07 15:00:00", stored)
}

func TestIncrementVisitorCount(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Visitor Count Fest", "Fest for visitor counts")

	t.Run("Increment Visitor Count", func(t *testing.T) {
		t.Parallel()
		initial, err := repo.ListVisitorCounts(t.Context(), festival.ID)
		assert.NoError(t, err)
		initialCount := 0
		for _, vc := range initial {
			initialCount += int(vc.Count)
		}

		err = repo.AddVisitorCount(t.Context(), festival.ID, time.Now(), 5)
		assert.NoError(t, err)

		after, err := repo.ListVisitorCounts(t.Context(), festival.ID)
		assert.NoError(t, err)
		afterCount := 0
		for _, vc := range after {
			afterCount += int(vc.Count)
		}

		assert.Equal(t, initialCount+5, afterCount)
	})
}

func TestDecrementVisitorCount(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Visitor Count Fest", "Fest for visitor counts")

	t.Run("Decrement Visitor Count", func(t *testing.T) {
		t.Parallel()
		err := repo.AddVisitorCount(t.Context(), festival.ID, time.Now(), 10)
		assert.NoError(t, err)

		initial, err := repo.ListVisitorCounts(t.Context(), festival.ID)
		assert.NoError(t, err)
		initialCount := 0
		for _, vc := range initial {
			initialCount += int(vc.Count)
		}

		err = repo.AddVisitorCount(t.Context(), festival.ID, time.Now(), -3)
		assert.NoError(t, err)

		after, err := repo.ListVisitorCounts(t.Context(), festival.ID)
		assert.NoError(t, err)
		afterCount := 0
		for _, vc := range after {
			afterCount += int(vc.Count)
		}

		assert.Equal(t, initialCount-3, afterCount)
	})
}

func TestDecrementVisitorCountBelowZero(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Visitor Count Fest", "Fest for visitor counts")

	t.Run("Decrement Visitor Count Below Zero", func(t *testing.T) {
		t.Parallel()
		err := repo.AddVisitorCount(t.Context(), festival.ID, time.Now(), 2)
		assert.NoError(t, err)

		err = repo.AddVisitorCount(t.Context(), festival.ID, time.Now(), -5)
		assert.NoError(t, err)

		after, err := repo.ListVisitorCounts(t.Context(), festival.ID)
		assert.NoError(t, err)
		afterCount := 0
		for _, vc := range after {
			afterCount += int(vc.Count)
		}

		assert.Equal(t, 0, afterCount)
	})
}
