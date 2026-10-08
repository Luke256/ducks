package gorm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
