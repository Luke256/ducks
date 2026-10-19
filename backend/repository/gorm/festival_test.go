package gorm

import (
	"testing"
	"time"

	"github.com/Luke256/ducks/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterFestival(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	date := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	festival, err := repo.RegisterFestival("Test Fest", "A fun festival", date)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, festival.ID)
	assert.True(t, festival.Date.Equal(date))
	stored, err := repo.GetFestivalByID(festival.ID)
	require.NoError(t, err)
	assert.True(t, stored.Date.Equal(date))
}

func TestGetFestivalByID(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Sample Fest", "Sample Description")

	t.Run("Get Existing Festival", func(t *testing.T) {
		t.Parallel()
		festival, err := repo.GetFestivalByID(festival.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Sample Fest", festival.Name)
		assert.Equal(t, "Sample Description", festival.Description)
	})

	t.Run("Get Non-Existent Festival", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		_, err := repo.GetFestivalByID(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("Get Zero UUID Festival", func(t *testing.T) {
		t.Parallel()
		zeroID := uuid.Nil
		_, err := repo.GetFestivalByID(zeroID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestGetAllFestivals(t *testing.T) {
	t.Parallel()
	repo := setup(t, s1)

	date := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	festival1, err := repo.RegisterFestival("Fest One", "First festival", date.AddDate(0, 0, 1))
	require.NoError(t, err)
	festival2, err := repo.RegisterFestival("Fest Two", "Second festival", date)
	require.NoError(t, err)
	festival3, err := repo.RegisterFestival("Fest Three", "Same date", festival1.Date)
	require.NoError(t, err)
	legacyID := uuid.New()
	require.NoError(t, repo.db.Exec("INSERT INTO festivals (id, name, description) VALUES (?, ?, ?)", legacyID, "Legacy festival", "No date").Error)

	festivals, err := repo.GetAllFestivals()
	require.NoError(t, err)
	var ids []uuid.UUID
	for _, fest := range festivals {
		if fest.ID == festival1.ID || fest.ID == festival2.ID || fest.ID == festival3.ID || fest.ID == legacyID {
			ids = append(ids, fest.ID)
		}
	}
	assert.Equal(t, []uuid.UUID{festival2.ID, festival1.ID, festival3.ID, legacyID}, ids)
}

func TestUpdateFestival(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Old Fest", "Old Description")

	t.Run("Update Festival", func(t *testing.T) {
		t.Parallel()
		date := festival.Date.AddDate(0, 0, 1)
		err := repo.UpdateFestival(festival.ID, "New Fest", "New Description", date)
		assert.NoError(t, err)
		updatedFestival, err := repo.GetFestivalByID(festival.ID)
		assert.NoError(t, err)
		assert.Equal(t, "New Fest", updatedFestival.Name)
		assert.Equal(t, "New Description", updatedFestival.Description)
		assert.True(t, updatedFestival.Date.Equal(date))
		assert.NoError(t, repo.UpdateFestival(festival.ID, updatedFestival.Name, updatedFestival.Description, date))
	})

	t.Run("Update Non-Existent Festival", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		err := repo.UpdateFestival(nonExistentID, "Ghost Fest", "No Description", festival.Date)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestDeleteFestival(t *testing.T) {
	t.Parallel()
	repo := setup(t, common)

	festival := mustCreateFestival(t, repo, "Delete Fest", "To be deleted")

	t.Run("Delete Existing Festival", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteFestival(festival.ID)
		assert.NoError(t, err)
		_, err = repo.GetFestivalByID(festival.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("Delete Non-Existent Festival", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		err := repo.DeleteFestival(nonExistentID)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}
