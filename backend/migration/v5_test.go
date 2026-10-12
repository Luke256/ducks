package migration

import (
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestV5MigratesExistingFestivalsUnlocked(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Festival{}))
	require.NoError(t, db.Migrator().DropColumn(&model.Festival{}, "Locked"))
	id := uuid.New()
	require.NoError(t, db.Exec(
		"INSERT INTO festivals (id, name, description) VALUES (?, ?, ?)",
		id, "Existing festival", "Existing description",
	).Error)

	require.NoError(t, v5().Migrate(db))
	var festival model.Festival
	require.NoError(t, db.First(&festival, "id = ?", id).Error)
	require.False(t, festival.Locked)
	require.Equal(t, "Existing festival", festival.Name)
	require.Equal(t, "Existing description", festival.Description)

	newID := uuid.New()
	require.NoError(t, db.Exec(
		"INSERT INTO festivals (id, name, description) VALUES (?, ?, ?)",
		newID, "New festival", "New description",
	).Error)
	var newFestival model.Festival
	require.NoError(t, db.First(&newFestival, "id = ?", newID).Error)
	require.False(t, newFestival.Locked)

	// 再実行しても、移行後にロックしたイベントを解除しない。
	require.NoError(t, db.Model(&festival).Update("locked", true).Error)
	require.NoError(t, v5().Migrate(db))
	require.NoError(t, db.First(&festival, "id = ?", id).Error)
	require.True(t, festival.Locked)
}
