package v1

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRegisterPoster(t *testing.T) {
	env := setup(t, common)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Poster Fest", "Festival for posters")

	resp := e.POST("/api/posters").
		WithMultipart().
		WithForm(map[string]any{
			"festival_id": fes.ID.String(),
			"name":        "Awesome Poster",
			"description": "This is an awesome poster.",
		}).
		WithFile("image", "poster_image.png", strings.NewReader("")).
		Expect().
		Status(201).
		JSON().
		Object()

	resp.Value("id").NotNull()
	resp.Value("festival_id").IsEqual(fes.ID.String())
	resp.Value("name").IsEqual("Awesome Poster")
	resp.Value("description").IsEqual("This is an awesome poster.")
	resp.Value("status").IsEqual(PosterStatusUncollected)
}

func TestListPostersByFestival(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "List Poster Fest", "Festival for listing posters")
	fes2 := env.mustCreateFestival(t, "Other Fest", "Another festival")

	poster1 := env.mustCreatePoster(t, fes.ID, "Poster One", "First poster")
	poster2 := env.mustCreatePoster(t, fes.ID, "Poster Two", "Second poster")

	_ = env.mustCreatePoster(t, fes2.ID, "Other Poster", "Poster in other festival")

	resp := e.GET("/api/festivals/{fesID}/posters", fes.ID.String()).
		Expect().
		Status(200).
		JSON().
		Object()

	array := resp.Value("posters").Array()
	array.Length().IsEqual(2)
	array.ContainsOnly(
		map[string]any{
			"id":          poster1.ID.String(),
			"festival_id": fes.ID.String(),
			"name": poster1.Name,
			"description": poster1.Description,
			"image_url":   poster1.ImageURL,
			"status":      poster1.Status,
		},
		map[string]any{
			"id":          poster2.ID.String(),
			"festival_id": fes.ID.String(),
			"name": poster2.Name,
			"description": poster2.Description,
			"image_url":   poster2.ImageURL,
			"status":      poster2.Status,
		},
	)
}

func TestGetPoster(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Get Poster Fest", "Festival for getting posters")
	poster := env.mustCreatePoster(t, fes.ID, "Gettable Poster", "Poster to be retrieved")

	t.Run("existing poster", func(t *testing.T) {
		resp := e.GET("/api/posters/%s", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival_id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual(poster.Name)
		resp.Value("description").IsEqual(poster.Description)
		resp.Value("image_url").IsEqual(poster.ImageURL)
		resp.Value("status").IsEqual(poster.Status)
	})

	t.Run("non-existent poster", func(t *testing.T) {
		nonExistentID := uuid.New()
		e.GET("/api/posters/%s", nonExistentID.String()).
			Expect().
			Status(404)
	})
}

func TestUpdatePoster(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Update Poster Fest", "Festival for updating posters")
	poster := env.mustCreatePoster(t, fes.ID, "Updatable Poster", "Poster to be updated")

	t.Run("update poster", func(t *testing.T) {
		e.PUT("/api/posters/%s", poster.ID.String()).
			WithJSON(map[string]any{
				"name": "Updated Poster Name",
				"description": "",
			}).
			Expect().
			Status(204)

		resp := e.GET("/api/posters/%s", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival_id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual("Updated Poster Name")
		resp.Value("description").IsEqual("")
		resp.Value("image_url").IsEqual(poster.ImageURL)
		resp.Value("status").IsEqual(poster.Status)
	})

	t.Run("update non-existent poster", func(t *testing.T) {
		e.PUT("/api/posters/00000000-0000-0000-0000-000000000000").
			WithJSON(map[string]any{
				"name": "Name",
				"description": "Description",
			}).
			Expect().
			Status(404)
	})
}

func TestUpdatePosterStatus(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fes := env.mustCreateFestival(t, "Status Poster Fest", "Festival for updating poster status")
	poster := env.mustCreatePoster(t, fes.ID, "Status Poster", "Poster to update status")

	t.Run("update poster status", func(t *testing.T) {
		e.PUT("/api/posters/%s/status", poster.ID.String()).
			WithJSON(map[string]any{
				"status": PosterStatusCollected,
			}).
			Expect().
			Status(204)

		resp := e.GET("/api/posters/%s", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival_id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual(poster.Name)
		resp.Value("description").IsEqual(poster.Description)
		resp.Value("image_url").IsEqual(poster.ImageURL)
		resp.Value("status").IsEqual(PosterStatusCollected)
	})

	t.Run("update non-existent poster status", func(t *testing.T) {
		e.PUT("/api/posters/00000000-0000-0000-0000-000000000000/status").
			WithJSON(map[string]any{
				"status": PosterStatusCollected,
			}).
			Expect().
			Status(404)
	})
}

func TestDeletePoster(t *testing.T) {
	env := setup(t, s1)
	e := env.R(t)

	fest := env.mustCreateFestival(t, "Delete Poster Fest", "Festival for deleting posters")
	poster := env.mustCreatePoster(t, fest.ID, "Deletable Poster", "Poster to be deleted")

	t.Run("delete existing poster", func(t *testing.T) {
		e.DELETE("/api/posters/%s", poster.ID.String()).
			Expect().
			Status(204)
		e.GET("/api/posters/%s", poster.ID.String()).
			Expect().
			Status(404)
	})

	t.Run("delete non-existent poster", func(t *testing.T) {
		e.DELETE("/api/posters/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}
