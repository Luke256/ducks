package v1

import (
	"testing"

	"github.com/google/uuid"
)

func TestIncrementVisitorCount(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	festival := env.mustCreateFestival(t, "Test Festival", "A festival for testing")

	t.Run("Increment Visitor Count - Success", func(t *testing.T) {
		e.POST("/api/visitors/{festival_id}", festival.ID).
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(204)
	})

	t.Run("Increment Visitor Count - Invalid Festival ID", func(t *testing.T) {
		e.POST("/api/visitors/{festival_id}", "invalid-uuid").
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(400)
	})

	t.Run("Increment Visitor Count - Festival Not Found", func(t *testing.T) {
		e.POST("/api/visitors/{festival_id}", uuid.New()).
			WithJSON(map[string]any{
				"amount": 5,
			}).
			Expect().
			Status(404)
	})
}

func TestGetVisitorCounts(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	festival := env.mustCreateFestival(t, "Test Festival", "A festival for testing")

	// Increment visitor count to have some data
	env.mustIncrementVisitorCount(t, festival.ID, 10)

	t.Run("Get Visitor Counts - Success", func(t *testing.T) {
		res := e.GET("/api/visitors/{festival_id}", festival.ID).
			Expect().
			Status(200).
			JSON().
			Object()

		res.Value("festival_id").IsEqual(festival.ID.String())
		res.Value("counts").Array().Length().Gt(0)
	})

	t.Run("Get Visitor Counts - Invalid Festival ID", func(t *testing.T) {
		e.GET("/api/visitors/{festival_id}", "invalid-uuid").
			Expect().
			Status(400)
	})

	t.Run("Get Visitor Counts - Festival Not Found", func(t *testing.T) {
		e.GET("/api/visitors/{festival_id}", uuid.New()).
			Expect().
			Status(200).JSON().Object().Value("counts").Array().Length().IsEqual(0)
	})
}
