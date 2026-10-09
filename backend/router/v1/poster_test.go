package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Luke256/ducks/repository"
	"github.com/Luke256/ducks/service/poster"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestPosterMultipartRejectsOversizedRequest(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/api/posters", strings.NewReader(""))
	req.ContentLength = int64(repository.MaxPosterImages*poster.MaxImageSize+(1<<20)) + 1
	c := echo.New().NewContext(req, httptest.NewRecorder())
	_, err := posterMultipart(c)
	var httpErr *echo.HTTPError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, http.StatusRequestEntityTooLarge, httpErr.Code)
}

func TestRegisterPoster(t *testing.T) {
	t.Parallel()
	env := setup(t, common)

	fes := env.mustCreateFestival(t, "Poster Fest", "Festival for posters")

	t.Run("register poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)

		resp := e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "Awesome Poster",
				"description": "This is an awesome poster.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			WithFile("image", "poster_image2.png", strings.NewReader("")).
			Expect().
			Status(201).
			JSON().
			Object()

		resp.Value("id").NotNull()
		resp.Value("festival").Object().Value("id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual("Awesome Poster")
		resp.Value("description").IsEqual("This is an awesome poster.")
		resp.Value("status").IsEqual(PosterStatusUncollected)
		resp.NotContainsKey("image_url")
		images := resp.Value("image").Array()
		images.Length().IsEqual(2)
		for i := range 2 {
			images.Element(i).Object().Value("id").String().NotEmpty()
			images.Element(i).Object().Value("url").String().NotEmpty()
		}
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "",
				"description": "No name poster.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(400)
	})

	t.Run("non-existent festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		nonExistentFesID := uuid.New()
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": nonExistentFesID.String(),
				"name":        "Ghost Poster",
				"description": "Poster for non-existent festival.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(404)
	})

	t.Run("missing image", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "Imageless Poster",
				"description": "Poster without an image.",
			}).
			Expect().
			Status(400)
	})

	t.Run("empty description", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "No Description Poster",
				"description": "",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(400)
	})

	t.Run("too long name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		longName := strings.Repeat("a", 65)
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        longName,
				"description": "Poster with too long name.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(400)
	})

	t.Run("duplicate poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "Duplicating Poster",
				"description": "This is a original poster.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(201)

		e.POST("/api/posters").
			WithMultipart().
			WithForm(map[string]any{
				"festival_id": fes.ID.String(),
				"name":        "Duplicating Poster",
				"description": "This is a duplicate poster.",
			}).
			WithFile("image", "poster_image.png", strings.NewReader("")).
			Expect().
			Status(409)
	})
}

func TestUpdatePosterImages(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)
	festival := env.mustCreateFestival(t, "Image edit API", "test")
	p := env.mustCreatePoster(t, festival.ID, "Editable images", "test")
	other := env.mustCreatePoster(t, festival.ID, "Foreign images", "test")
	path := "/api/posters/{id}/images"

	// 同じポスターの編集結果に依存するサブテストは順番に実行する。
	t.Run("append multiple images", func(t *testing.T) {
		e := env.R(t)
		resp := e.PATCH(path, p.ID.String()).WithMultipart().
			WithFile("image", "a.png", strings.NewReader("")).
			WithFile("image", "b.png", strings.NewReader("")).
			Expect().Status(200).JSON().Object()
		resp.NotContainsKey("image_url")
		resp.Value("image").Array().Length().IsEqual(3)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		require.Contains(t, stored.Images, p.Images[0])
	})
	t.Run("replace one and retain others", func(t *testing.T) {
		e := env.R(t)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFormField("delete_image_ids", p.Images[0].ID).
			WithFile("image", "replacement.png", strings.NewReader("")).
			Expect().Status(200).JSON().Object().Value("image").Array().Length().IsEqual(3)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		require.NotContains(t, stored.Images, p.Images[0])
	})
	t.Run("reject foreign image", func(t *testing.T) {
		e := env.R(t)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFormField("delete_image_ids", other.Images[0].ID).
			Expect().Status(400)
		stored, err := env.PM.Get(other.ID)
		require.NoError(t, err)
		require.Equal(t, other.Images, stored.Images)
	})
	t.Run("reject missing and duplicate deletion", func(t *testing.T) {
		e := env.R(t)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFormField("delete_image_ids", uuid.NewString()).Expect().Status(400)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFormField("delete_image_ids", stored.Images[0].ID).
			WithFormField("delete_image_ids", stored.Images[0].ID).Expect().Status(400)
	})
	t.Run("reject deleting all images without uploads", func(t *testing.T) {
		e := env.R(t)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		req := e.PATCH(path, p.ID.String()).WithMultipart()
		for _, image := range stored.Images {
			req.WithFormField("delete_image_ids", image.ID)
		}
		req.Expect().Status(400)
		e.GET("/api/posters/{id}", p.ID.String()).Expect().Status(200).
			JSON().Object().Value("image").IsEqual(stored.Images)
	})
	t.Run("replace all images with one new image", func(t *testing.T) {
		e := env.R(t)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		req := e.PATCH(path, p.ID.String()).WithMultipart().
			WithFile("image", "replacement.png", strings.NewReader(""))
		for _, image := range stored.Images {
			req.WithFormField("delete_image_ids", image.ID)
		}
		req.Expect().Status(200).JSON().Object().Value("image").Array().Length().IsEqual(1)
		updated, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		require.Len(t, updated.Images, 1)
		require.NotContains(t, stored.Images, updated.Images[0])
	})
	t.Run("reject deleting the last image", func(t *testing.T) {
		e := env.R(t)
		stored, err := env.PM.Get(p.ID)
		require.NoError(t, err)
		require.Len(t, stored.Images, 1)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFormField("delete_image_ids", stored.Images[0].ID).Expect().Status(400)
		e.GET("/api/posters/{id}", p.ID.String()).Expect().Status(200).
			JSON().Object().Value("image").IsEqual(stored.Images)
	})
	t.Run("missing poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH(path, uuid.NewString()).WithMultipart().WithFile("image", "new.png", strings.NewReader("")).Expect().Status(404)
		e.PATCH(path, "invalid").WithMultipart().WithFile("image", "new.png", strings.NewReader("")).Expect().Status(404)
	})
	t.Run("empty changes", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH(path, p.ID.String()).WithMultipart().WithFormField("unused", "value").Expect().Status(400)
	})
	t.Run("requires multipart", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH(path, p.ID.String()).WithJSON(map[string]any{"delete_image_ids": []string{other.Images[0].ID}}).Expect().Status(400)
	})
	t.Run("image count limit", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		req := e.PATCH(path, p.ID.String()).WithMultipart()
		for range repository.MaxPosterImages + 1 {
			req.WithFile("image", "new.png", strings.NewReader(""))
		}
		req.Expect().Status(400)
	})
	t.Run("individual file size limit", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH(path, p.ID.String()).WithMultipart().
			WithFile("image", "large.png", strings.NewReader(strings.Repeat("a", poster.MaxImageSize+1))).
			Expect().Status(413)
	})
}

func TestListPostersByFestival(t *testing.T) {
	t.Parallel()
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
			"id": poster1.ID.String(),
			"festival": map[string]any{
				"id":          fes.ID.String(),
				"name":        fes.Name,
				"description": fes.Description,
			},
			"name":        poster1.Name,
			"description": poster1.Description,
			"image":       poster1.Images,
			"status":      poster1.Status,
		},
		map[string]any{
			"id": poster2.ID.String(),
			"festival": map[string]any{
				"id":          fes.ID.String(),
				"name":        fes.Name,
				"description": fes.Description,
			},
			"name":        poster2.Name,
			"description": poster2.Description,
			"image":       poster2.Images,
			"status":      poster2.Status,
		},
	)
}

func TestGetPoster(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)

	fes := env.mustCreateFestival(t, "Get Poster Fest", "Festival for getting posters")
	poster := env.mustCreatePoster(t, fes.ID, "Gettable Poster", "Poster to be retrieved")

	t.Run("existing poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		resp := e.GET("/api/posters/{posterID}", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival").Object().Value("id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual(poster.Name)
		resp.Value("description").IsEqual(poster.Description)
		resp.Value("image").IsEqual(poster.Images)
		resp.Value("status").IsEqual(poster.Status)
	})

	t.Run("non-existent poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		nonExistentID := uuid.New()
		e.GET("/api/posters/{posterID}", nonExistentID.String()).
			Expect().
			Status(404)
	})
}

func TestGetPosterByFestivalAndName(t *testing.T) {
	t.Parallel()
	env := setup(t, common)

	fes := env.mustCreateFestival(t, "Name Poster Fest", "Festival for getting posters by name")
	poster := env.mustCreatePoster(t, fes.ID, "Unique Poster", "Poster with unique name")

	t.Run("existing poster by name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		resp := e.GET("/api/posters/{festivalID}/{posterName}", fes.ID.String(), poster.Name).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival").Object().Value("id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual(poster.Name)
		resp.Value("description").IsEqual(poster.Description)
		resp.Value("image").IsEqual(poster.Images)
		resp.Value("status").IsEqual(poster.Status)
	})

	t.Run("non-existent poster by name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.GET("/api/posters/{festivalID}/{posterName}", fes.ID.String(), "NonExistentPoster").
			Expect().
			Status(404)
	})

	t.Run("non-existent festival", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		nonExistentFesID := uuid.New()
		e.GET("/api/posters/{festivalID}/{posterName}", nonExistentFesID.String(), poster.Name).
			Expect().
			Status(404)
	})
}

func TestUpdatePoster(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)

	fes := env.mustCreateFestival(t, "Update Poster Fest", "Festival for updating posters")
	poster := env.mustCreatePoster(t, fes.ID, "Updatable Poster", "Poster to be updated")

	t.Run("update poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PUT("/api/posters/{posterID}", poster.ID.String()).
			WithJSON(map[string]any{
				"name":        "Updated Poster Name",
				"description": "Updated description.",
			}).
			Expect().
			Status(204)

		resp := e.GET("/api/posters/{posterID}", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival").Object().Value("id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual("Updated Poster Name")
		resp.Value("description").IsEqual("Updated description.")
		resp.Value("image").IsEqual(poster.Images)
		resp.Value("status").IsEqual(poster.Status)
	})

	t.Run("update non-existent poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PUT("/api/posters/00000000-0000-0000-0000-000000000000").
			WithJSON(map[string]any{
				"name":        "Name",
				"description": "Description",
			}).
			Expect().
			Status(404)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PUT("/api/posters/{posterID}", poster.ID.String()).
			WithJSON(map[string]any{
				"name":        "",
				"description": "Description",
			}).
			Expect().
			Status(400)
	})

	t.Run("too long name", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		longName := strings.Repeat("b", 65)
		e.PUT("/api/posters/{posterID}", poster.ID.String()).
			WithJSON(map[string]any{
				"name":        longName,
				"description": "Description",
			}).
			Expect().
			Status(400)
	})

	t.Run("empty description", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PUT("/api/posters/{posterID}", poster.ID.String()).
			WithJSON(map[string]any{
				"name":        "Name",
				"description": "",
			}).
			Expect().
			Status(400)
	})
}

func TestUpdatePosterStatus(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)

	fes := env.mustCreateFestival(t, "Status Poster Fest", "Festival for updating poster status")
	poster := env.mustCreatePoster(t, fes.ID, "Status Poster", "Poster to update status")

	t.Run("update poster status", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH("/api/posters/{posterID}/status", poster.ID.String()).
			WithJSON(map[string]any{
				"status": PosterStatusCollected,
			}).
			Expect().
			Status(204)

		resp := e.GET("/api/posters/{posterID}", poster.ID.String()).
			Expect().
			Status(200).
			JSON().
			Object()
		resp.Value("id").IsEqual(poster.ID.String())
		resp.Value("festival").Object().Value("id").IsEqual(fes.ID.String())
		resp.Value("name").IsEqual(poster.Name)
		resp.Value("description").IsEqual(poster.Description)
		resp.Value("image").IsEqual(poster.Images)
		resp.Value("status").IsEqual(PosterStatusCollected)
	})

	t.Run("update non-existent poster status", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.PATCH("/api/posters/00000000-0000-0000-0000-000000000000/status").
			WithJSON(map[string]any{
				"status": PosterStatusCollected,
			}).
			Expect().
			Status(404)
	})
}

func TestDeletePoster(t *testing.T) {
	t.Parallel()
	env := setup(t, s1)

	fest := env.mustCreateFestival(t, "Delete Poster Fest", "Festival for deleting posters")
	poster := env.mustCreatePoster(t, fest.ID, "Deletable Poster", "Poster to be deleted")

	t.Run("delete existing poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.DELETE("/api/posters/{posterID}", poster.ID.String()).
			Expect().
			Status(204)
		e.GET("/api/posters/{posterID}", poster.ID.String()).
			Expect().
			Status(404)
	})

	t.Run("delete non-existent poster", func(t *testing.T) {
		t.Parallel()
		e := env.R(t)
		e.DELETE("/api/posters/00000000-0000-0000-0000-000000000000").
			Expect().
			Status(404)
	})
}
