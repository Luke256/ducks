package v1

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateFestival(t *testing.T) {
	t.Parallel()
	env := setup(t, common)
	e := env.R(t)

	resp := e.POST("/api/festivals").
		WithJSON(map[string]any{
			"name":        "Test Festival",
			"description": "This is a test festival.",
			"date":        "2026-10-18",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	resp.Value("id").NotNull()
	resp.Value("name").IsEqual("Test Festival")
	resp.Value("description").IsEqual("This is a test festival.")
	resp.Value("locked").IsEqual(false)
	resp.Value("date").IsEqual("2026-10-18T00:00:00Z")
}

func TestListFestivals(t *testing.T) {
	t.Parallel()
	// 一覧の件数を検証するため、他のテストが祭りを追加しないDBを使う。
	env := setup(t, s4)
	e := env.R(t)

	fest1, err := env.FM.Create("Festival 1", "Description 1", time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	fest2, err := env.FM.Create("Festival 2", "Description 2", time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	resp := e.GET("/api/festivals").
		Expect().
		Status(200).
		JSON().
		Object()

	array := resp.Value("festivals").Array()
	array.Length().IsEqual(2)
	array.ContainsOnly(
		map[string]any{
			"id":          fest1.ID.String(),
			"name":        fest1.Name,
			"description": fest1.Description,
			"date":        fest1.Date.Format(time.RFC3339),
			"locked":      false,
		},
		map[string]any{
			"id":          fest2.ID.String(),
			"name":        fest2.Name,
			"description": fest2.Description,
			"date":        fest2.Date.Format(time.RFC3339),
			"locked":      false,
		},
	)
	array.Value(0).Object().Value("id").IsEqual(fest2.ID.String())
	array.Value(1).Object().Value("id").IsEqual(fest1.ID.String())
}

func TestGetFestival(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)
	fest := env.mustCreateFestival(t, "Festival", "Description")

	t.Run("existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		resp := e.GET(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(fest.ID.String())
		resp.Value("name").IsEqual(fest.Name)
		resp.Value("description").IsEqual(fest.Description)
		resp.Value("date").IsEqual(fest.Date.Format(time.RFC3339))
	})

	t.Run("non-existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.GET("/api/festivals/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}

func TestEditFestival(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)
	fest := env.mustCreateFestival(t, "Old Name", "Old Description")

	t.Run("edit existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		resp := e.PUT(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			WithJSON(map[string]any{
				"name":        "New Name",
				"description": "",
				"date":        "2026-10-19",
			}).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(fest.ID.String())
		resp.Value("name").IsEqual("New Name")
		resp.Value("description").IsEqual("")
		resp.Value("date").IsEqual("2026-10-19T00:00:00Z")
		e.GET(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).Expect().Status(200).
			JSON().Object().Value("date").IsEqual("2026-10-19T00:00:00Z")

		e.PUT(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			WithJSON(map[string]any{
				"name":        "New Name",
				"description": "",
				"date":        "2026-10-19",
			}).
			Expect().
			Status(200)
	})

	t.Run("edit non-existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PUT("/api/festivals/00000000-0000-0000-0000-000000000000").
			WithJSON(map[string]any{
				"name":        "Name",
				"description": "Description",
				"date":        "2026-10-19",
			}).
			Expect().
			Status(404)
	})
}

func TestDeleteFestival(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)
	fest := env.mustCreateFestival(t, "To Be Deleted", "Description")

	t.Run("delete existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.DELETE(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(204)
		e.GET(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(404)
	})

	t.Run("delete non-existing festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.DELETE("/api/festivals/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}
