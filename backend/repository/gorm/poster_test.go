package gorm

import (
	"testing"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPoster(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")

	t.Run("Register Poster", func(t *testing.T) {
		poster, err := repo.RegisterPoster(festival.ID, "PosterReg", "desc", "img-1")
		assert.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, poster.ID)
	})
}

func TestGetPostersByFestivalID(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	posterID1 := mustCreatePoster(t, repo, festival.ID, "PosterOne", "desc1", "img-1")
	posterID2 := mustCreatePoster(t, repo, festival.ID, "PosterTwo", "desc2", "img-2")

	t.Run("Get Posters by FestivalID", func(t *testing.T) {
		posters, err := repo.GetPostersByFestivalID(festival.ID)
		assert.NoError(t, err)
		assert.Len(t, posters, 2)
		foundPoster1 := false
		foundPoster2 := false
		for _, p := range posters {
			if p.ID == posterID1.ID {
				foundPoster1 = true
			}
			if p.ID == posterID2.ID {
				foundPoster2 = true
			}
		}
		assert.True(t, foundPoster1, "PosterOne not found in GetPostersByFestivalID result")
		assert.True(t, foundPoster2, "PosterTwo not found in GetPostersByFestivalID result")
	})
}

func TestGetPosterByID(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	
	poster := mustCreatePoster(t, repo, festival.ID, "PosterQuery", "desc", "img-2")

	t.Run("Get Existing Poster", func(t *testing.T) {
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "PosterQuery", p.PosterName)
		assert.Equal(t, "desc", p.Description)
		assert.Equal(t, "img-2", p.ImageID)
		assert.Equal(t, festival.ID, p.FestivalID)
	})

	t.Run("Get Non-Existent Poster", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := repo.GetPosterByID(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetPosterByFestivalIDAndPosterName(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")

	poster := mustCreatePoster(t, repo, festival.ID, "PosterByName", "desc-name", "img-name")

	t.Run("Get Existing Poster by FestivalID and PosterName", func(t *testing.T) {
		p, err := repo.GetPosterByFestivalIDAndPosterName(festival.ID, "PosterByName")
		assert.NoError(t, err)
		assert.Equal(t, poster.ID, p.ID)
		assert.Equal(t, "desc-name", p.Description)
		assert.Equal(t, "img-name", p.ImageID)
	})
}

func TestUpdatePoster(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterToUpdate", "old-desc", "old-img")

	t.Run("Update Poster Info", func(t *testing.T) {
		err := repo.UpdatePoster(poster.ID, "UpdatedPoster", "new-desc")
		assert.NoError(t, err)
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedPoster", p.PosterName)
		assert.Equal(t, "new-desc", p.Description)
	})

	t.Run("Update Non-Existent Poster", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.UpdatePoster(nonExistentID, "NoPoster", "no-desc")
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestUpdatePosterStatus(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterStatus", "status-desc", "status-img")

	t.Run("Update Poster Status", func(t *testing.T) {
		err := repo.UpdatePosterStatus(poster.ID, "collected")
		assert.NoError(t, err)
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "collected", p.Status)
	})

	t.Run("Update Non-Existent Poster Status", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.UpdatePosterStatus(nonExistentID, "lost")
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeletePoster(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterToDelete", "del-desc", "del-img")

	t.Run("Delete Existing Poster", func(t *testing.T) {
		err := repo.DeletePoster(poster.ID)
		assert.NoError(t, err)
		_, err = repo.GetPosterByID(poster.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Poster", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.DeletePoster(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}