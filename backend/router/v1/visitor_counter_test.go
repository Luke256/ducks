package v1

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIncrementVisitorCount(t *testing.T) {
	t.Parallel()
	env := setup(t, common)

	festival := env.mustCreateFestival(t, "Test Festival", "A festival for testing")

	t.Run("Increment Visitor Count - Success", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/visitors/{festival_id}", festival.ID).
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(204)
	})

	t.Run("Increment Visitor Count - Invalid Festival ID", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/visitors/{festival_id}", "invalid-uuid").
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(400)
	})

	t.Run("Increment Visitor Count - Festival Not Found", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/visitors/{festival_id}", uuid.New()).
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(404)
	})
}

func TestGetVisitorCounts(t *testing.T) {
	t.Parallel()
	env := setup(t, common)

	festival := env.mustCreateFestival(t, "Test Festival", "A festival for testing")

	// Increment visitor count to have some data
	timestamp := time.Date(2026, 10, 7, 15, 9, 59, 0, time.UTC)
	require.NoError(t, env.Repo.AddVisitorCount(t.Context(), festival.ID, timestamp, 10))

	t.Run("Get Visitor Counts - Success", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		res := e.GET("/api/visitors/{festival_id}", festival.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("festival_id").IsEqual(festival.ID.String())
		counts := res.Value("counts").Array()
		counts.Length().IsEqual(1)
		counts.Value(0).Object().Value("bucket_start").IsEqual("2026-10-08T00:00:00+09:00")
		counts.Value(0).Object().Value("count").IsEqual(10)
		stored, err := env.Repo.ListVisitorCounts(t.Context(), festival.ID)
		require.NoError(t, err)
		require.Len(t, stored, 1)
		require.Equal(t, time.Date(2026, 10, 7, 15, 0, 0, 0, time.UTC), stored[0].BucketStart)
	})

	t.Run("Get Visitor Counts - Invalid Festival ID", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.GET("/api/visitors/{festival_id}", "invalid-uuid").
			Expect().
			Status(400)
	})

	t.Run("Get Visitor Counts - Festival Not Found", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.GET("/api/visitors/{festival_id}", uuid.New()).
			Expect().
			Status(200).JSON().Object().Value("counts").Array().Length().IsEqual(0)
	})
}
