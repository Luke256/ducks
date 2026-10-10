package gorm

import (
	"sync"
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterPoster(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")

	t.Run("Register Poster", func(t *testing.T) {
		t.Parallel()
		poster, err := repo.RegisterPoster(festival.ID, "PosterReg", "desc", []string{"img-reg-1", "img-reg-2"})
		assert.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, poster.ID)
		stored, err := repo.GetPosterByID(poster.ID)
		require.NoError(t, err)
		assert.ElementsMatch(t, poster.Images, stored.Images)
		assert.Len(t, stored.Images, 2)
	})
}

func TestGetPostersByFestivalID(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	posterID1 := mustCreatePoster(t, repo, festival.ID, "PosterOne", "desc1", "img-list-1")
	posterID2 := mustCreatePoster(t, repo, festival.ID, "PosterTwo", "desc2", "img-list-2")

	t.Run("Get Posters by FestivalID", func(t *testing.T) {
		t.Parallel()
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
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")

	poster := mustCreatePoster(t, repo, festival.ID, "PosterQuery", "desc", "img-2")

	t.Run("Get Existing Poster", func(t *testing.T) {
		t.Parallel()
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "PosterQuery", p.PosterName)
		assert.Equal(t, "desc", p.Description)
		require.Len(t, p.Images, 1)
		assert.Equal(t, "img-2", p.Images[0].ID)
		assert.Equal(t, festival.ID, p.FestivalID)
		assert.Equal(t, festival.ID, p.Festival.ID)
	})

	t.Run("Get Non-Existent Poster", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		_, err := repo.GetPosterByID(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetPosterByFestivalIDAndPosterName(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")

	poster := mustCreatePoster(t, repo, festival.ID, "PosterByName", "desc-name", "img-name")

	t.Run("Get Existing Poster by FestivalID and PosterName", func(t *testing.T) {
		t.Parallel()
		p, err := repo.GetPosterByFestivalIDAndPosterName(festival.ID, "PosterByName")
		assert.NoError(t, err)
		assert.Equal(t, poster.ID, p.ID)
		assert.Equal(t, "desc-name", p.Description)
		require.Len(t, p.Images, 1)
		assert.Equal(t, "img-name", p.Images[0].ID)
		assert.Equal(t, festival.ID, p.Festival.ID)
	})
}

func TestUpdatePoster(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterToUpdate", "old-desc", "old-img")

	t.Run("Update Poster Info", func(t *testing.T) {
		t.Parallel()
		err := repo.UpdatePoster(poster.ID, "UpdatedPoster", "new-desc")
		assert.NoError(t, err)
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedPoster", p.PosterName)
		assert.Equal(t, "new-desc", p.Description)
		assert.NoError(t, repo.UpdatePoster(poster.ID, p.PosterName, p.Description))
	})

	t.Run("Update Non-Existent Poster", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		err := repo.UpdatePoster(nonExistentID, "NoPoster", "no-desc")
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestUpdatePosterStatus(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterStatus", "status-desc", "status-img")

	t.Run("Update Poster Status", func(t *testing.T) {
		t.Parallel()
		err := repo.UpdatePosterStatus(poster.ID, "collected")
		assert.NoError(t, err)
		p, err := repo.GetPosterByID(poster.ID)
		assert.NoError(t, err)
		assert.Equal(t, "collected", p.Status)
		assert.NoError(t, repo.UpdatePosterStatus(poster.ID, p.Status))
	})

	t.Run("Update Non-Existent Poster Status", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		err := repo.UpdatePosterStatus(nonExistentID, "lost")
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestRegisterPosterDoesNotReassignImages(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)
	festival := mustCreateFestival(t, repo, "Image owner", "test")
	imageID := uuid.NewString()
	owner := mustCreatePoster(t, repo, festival.ID, "Owner", "test", imageID)
	_, err := repo.RegisterPoster(festival.ID, "Collision", "test", []string{imageID})
	require.Error(t, err)
	stored, err := repo.GetPosterByID(owner.ID)
	require.NoError(t, err)
	require.Equal(t, owner.Images, stored.Images)
	posters, err := repo.GetPostersByFestivalID(festival.ID)
	require.NoError(t, err)
	require.Len(t, posters, 1)
}

func TestUpdatePosterImages(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)
	festival := mustCreateFestival(t, repo, "Image editing", "test")
	oldID, keptID, foreignID := uuid.NewString(), uuid.NewString(), uuid.NewString()
	poster, err := repo.RegisterPoster(festival.ID, "Editable", "test", []string{oldID, keptID})
	require.NoError(t, err)
	other := mustCreatePoster(t, repo, festival.ID, "Other", "test", foreignID)

	for _, deleteIDs := range [][]string{{foreignID}, {oldID, oldID}} {
		_, err := repo.UpdatePosterImages(poster.ID, nil, deleteIDs)
		require.ErrorIs(t, err, repository.ErrInvalidPosterImages)
	}
	// 削除後のINSERTが失敗しても、旧画像の削除はロールバックされる。
	_, err = repo.UpdatePosterImages(poster.ID, []string{foreignID}, []string{oldID})
	require.Error(t, err)
	stored, err := repo.GetPosterByID(poster.ID)
	require.NoError(t, err)
	require.ElementsMatch(t, poster.Images, stored.Images)
	foreign, err := repo.GetPosterByID(other.ID)
	require.NoError(t, err)
	require.Equal(t, other.Images, foreign.Images)

	newID := uuid.NewString()
	images, err := repo.UpdatePosterImages(poster.ID, []string{newID}, []string{oldID})
	require.NoError(t, err)
	require.ElementsMatch(t, []model.PosterImage{{ID: keptID, PosterID: poster.ID}, {ID: newID, PosterID: poster.ID}}, images)
	_, err = repo.UpdatePosterImages(poster.ID, nil, []string{keptID, newID})
	require.ErrorIs(t, err, repository.ErrInvalidPosterImages)
	stored, err = repo.GetPosterByID(poster.ID)
	require.NoError(t, err)
	require.ElementsMatch(t, images, stored.Images)
	replacementID := uuid.NewString()
	images, err = repo.UpdatePosterImages(poster.ID, []string{replacementID}, []string{keptID, newID})
	require.NoError(t, err)
	require.Equal(t, []model.PosterImage{{ID: replacementID, PosterID: poster.ID}}, images)
	_, err = repo.UpdatePosterImages(uuid.New(), []string{uuid.NewString()}, nil)
	require.ErrorIs(t, err, repository.ErrNotFound)
}

func TestConcurrentPosterImageChangesRespectLimit(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)
	festival := mustCreateFestival(t, repo, "Image limit", "test")
	ids := make([]string, repository.MaxPosterImages-1)
	for i := range ids {
		ids[i] = uuid.NewString()
	}
	poster, err := repo.RegisterPoster(festival.ID, "Nearly full", "test", ids)
	require.NoError(t, err)
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	start := make(chan struct{})
	for range 2 {
		wg.Go(func() {
			<-start
			_, err := repo.UpdatePosterImages(poster.ID, []string{uuid.NewString()}, nil)
			errors <- err
		})
	}
	close(start)
	wg.Wait()
	close(errors)
	var successes int
	for err := range errors {
		if err == nil {
			successes++
		} else {
			require.ErrorIs(t, err, repository.ErrInvalidPosterImages)
		}
	}
	require.Equal(t, 1, successes)
	stored, err := repo.GetPosterByID(poster.ID)
	require.NoError(t, err)
	require.Len(t, stored.Images, repository.MaxPosterImages)
}

func TestDeletePoster(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Poster Fest", "Fest for posters")
	poster := mustCreatePoster(t, repo, festival.ID, "PosterToDelete", "del-desc", "del-img")

	t.Run("Delete Existing Poster", func(t *testing.T) {
		t.Parallel()
		err := repo.DeletePoster(poster.ID)
		assert.NoError(t, err)
		_, err = repo.GetPosterByID(poster.ID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Delete Non-Existent Poster", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		err := repo.DeletePoster(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
