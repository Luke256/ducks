package gorm

import (
	"testing"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterFestival(t *testing.T) {
	repo := setup(t, common)

	id, err := repo.RegisterFestival("Test Fest", "A fun festival")
	assert.NoError(t, err)
	assert.NotEqual(t, "", id)
}

func TestGetFestivalByID(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Sample Fest", "Sample Description")

	t.Run("Get Existing Festival", func(t *testing.T) {
		festival, err := repo.GetFestivalByID(festival.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Sample Fest", festival.Name)
		assert.Equal(t, "Sample Description", festival.Description)
	})

	t.Run("Get Non-Existent Festival", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := repo.GetFestivalByID(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Get Zero UUID Festival", func(t *testing.T) {
		zeroID := uuid.Nil
		_, err := repo.GetFestivalByID(zeroID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetAllFestivals(t *testing.T) {
	repo := setup(t, s1)

	festival1 := mustCreateFestival(t, repo, "Fest One", "First festival")
	festival2 := mustCreateFestival(t, repo, "Fest Two", "Second festival")

	festivals, err := repo.GetAllFestivals()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(festivals), 2)
	foundFest1 := false
	foundFest2 := false
	for _, fest := range festivals {
		if fest.ID == festival1.ID {
			foundFest1 = true
		}
		if fest.ID == festival2.ID {
			foundFest2 = true
		}
	}
	assert.True(t, foundFest1, "Festival One not found in GetAllFestivals result")
	assert.True(t, foundFest2, "Festival Two not found in GetAllFestivals result")
}

func TestUpdateFestival(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Old Fest", "Old Description")

	t.Run("Update Festival", func(t *testing.T) {
		err := repo.UpdateFestival(festival.ID, "New Fest", "New Description")
		assert.NoError(t, err)
		updatedFestival, err := repo.GetFestivalByID(festival.ID)
		assert.NoError(t, err)
		assert.Equal(t, "New Fest", updatedFestival.Name)
		assert.Equal(t, "New Description", updatedFestival.Description)
	})

	t.Run("Update Non-Existent Festival", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.UpdateFestival(nonExistentID, "Ghost Fest", "No Description")
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteFestival(t *testing.T) {
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Delete Fest", "To be deleted")

	t.Run("Delete Existing Festival", func(t *testing.T) {
		err := repo.DeleteFestival(festival.ID)
		assert.NoError(t, err)
		_, err = repo.GetFestivalByID(festival.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("Delete Non-Existent Festival", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.DeleteFestival(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}