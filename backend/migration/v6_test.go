package migration

import (
	"testing"
	"time"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestV6MigratesExistingFestivals(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Festival{}))
	require.NoError(t, db.Migrator().DropColumn(&model.Festival{}, "Date"))
	id := uuid.New()
	require.NoError(t, db.Exec(
		"INSERT INTO festivals (id, name, description, locked) VALUES (?, ?, ?, ?)",
		id, "Existing festival", "Existing description", true,
	).Error)

	require.NoError(t, v6().Migrate(db))
	require.True(t, db.Migrator().HasColumn(&model.Festival{}, "Date"))
	var festival model.Festival
	require.NoError(t, db.First(&festival, "id = ?", id).Error)
	require.Equal(t, "Existing festival", festival.Name)
	require.Equal(t, "Existing description", festival.Description)
	require.True(t, festival.Locked)
	require.True(t, festival.Date.IsZero())

	date := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	require.NoError(t, db.Model(&festival).Update("date", date).Error)
	require.NoError(t, v6().Migrate(db))
	require.NoError(t, db.First(&festival, "id = ?", id).Error)
	require.True(t, festival.Date.Equal(date))
}
