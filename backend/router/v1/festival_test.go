package v1

import (
	"fmt"
	"testing"
)

func TestCreateFestival(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	resp := e.POST("/api/festivals").
		WithJSON(map[string]interface{}{
			"name": "Test Festival",
			"description": "This is a test festival.",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	
	resp.Value("id").NotNull()
	resp.Value("name").IsEqual("Test Festival")
	resp.Value("description").IsEqual("This is a test festival.")
}

func TestListFestivals(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fest1 := env.mustCreateFestival(t, "Festival 1", "Description 1")
	fest2 := env.mustCreateFestival(t, "Festival 2", "Description 2")

	resp := e.GET("/api/festivals").
		Expect().
		Status(200).
		JSON().
		Object()

	array := resp.Value("festivals").Array()
	array.Length().IsEqual(2)
	array.ContainsOnly(
		map[string]any{
			"id": fest1.ID.String(),
			"name": fest1.Name,
			"description": fest1.Description,
		},
		map[string]any{
			"id": fest2.ID.String(),
			"name": fest2.Name,
			"description": fest2.Description,
		},
	)
}

func TestGetFestival(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)
	fest := env.mustCreateFestival(t, "Festival", "Description")

	t.Run("existing festival", func(t *testing.T) {
		resp := e.GET(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(fest.ID.String())
		resp.Value("name").IsEqual(fest.Name)
		resp.Value("description").IsEqual(fest.Description)
	})

	t.Run("non-existing festival", func(t *testing.T) {
		e.GET("/api/festivals/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}

func TestEditFestival(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)
	fest := env.mustCreateFestival(t, "Old Name", "Old Description")

	t.Run("edit existing festival", func(t *testing.T) {
		resp := e.PUT(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			WithJSON(map[string]interface{}{
				"name": "New Name",
				"description": "",
			}).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(fest.ID.String())
		resp.Value("name").IsEqual("New Name")
		resp.Value("description").IsEqual("")
	})

	t.Run("edit non-existing festival", func(t *testing.T) {
		e.PUT("/api/festivals/00000000-0000-0000-0000-000000000000").
			WithJSON(map[string]any{
				"name": "Name",
				"description": "Description",
			}).
			Expect().
			Status(404)
	})
}

func TestDeleteFestival(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)
	fest := env.mustCreateFestival(t, "To Be Deleted", "Description")

	t.Run("delete existing festival", func(t *testing.T) {
		e.DELETE(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(204)
		e.GET(fmt.Sprintf("/api/festivals/%s", fest.ID.String())).
			Expect().
			Status(404)
	})

	t.Run("delete non-existing festival", func(t *testing.T) {
		e.DELETE("/api/festivals/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}