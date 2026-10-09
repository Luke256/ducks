package migration

import (
	"fmt"
	"testing"

	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/utils"
	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func migrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	config := &driverMysql.Config{
		User:                 utils.GetEnvOrDefault("NS_MARIADB_USERNAME", "root"),
		Passwd:               utils.GetEnvOrDefault("NS_MARIADB_PASSWORD", "password"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", utils.GetEnvOrDefault("NS_MARIADB_HOST", "localhost"), utils.GetEnvOrDefault("NS_MARIADB_PORT", "3307")),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	admin, err := gorm.Open(mysql.Open(config.FormatDSN()), &gorm.Config{})
	require.NoError(t, err)
	adminSQL, err := admin.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = adminSQL.Close() })

	name := "ducks-migration-test-" + uuid.NewString()
	require.NoError(t, admin.Exec("CREATE DATABASE `"+name+"`").Error)
	t.Cleanup(func() { require.NoError(t, admin.Exec("DROP DATABASE `"+name+"`").Error) })
	config.DBName = name
	db, err := gorm.Open(mysql.Open(config.FormatDSN()), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

func TestV3MigratesExistingImages(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Poster{}, &model.PosterImage{}))
	require.NoError(t, db.Exec("ALTER TABLE posters ADD COLUMN image_id TEXT NOT NULL").Error)
	festival := model.Festival{ID: uuid.New(), Name: "test"}
	require.NoError(t, db.Create(&festival).Error)
	poster := model.Poster{
		ID: uuid.New(), FestivalID: festival.ID, PosterName: "test",
		Description: "location", Status: "uncollected",
	}
	imageID := uuid.NewString() + ".webp"
	require.NoError(t, db.Exec(
		"INSERT INTO posters (id, festival_id, poster_name, description, status, image_id) VALUES (?, ?, ?, ?, ?, ?)",
		poster.ID, poster.FestivalID, poster.PosterName, poster.Description, poster.Status, imageID,
	).Error)
	// 画像のコピーだけが完了した状態からでも再実行できる。
	require.NoError(t, db.Create(&model.PosterImage{ID: imageID, PosterID: poster.ID}).Error)

	require.NoError(t, v3().Migrate(db))
	require.False(t, db.Migrator().HasColumn(&model.Poster{}, "image_id"))
	require.NoError(t, db.Preload("Images").First(&poster, "id = ?", poster.ID).Error)
	require.Equal(t, []model.PosterImage{{ID: imageID, PosterID: poster.ID}}, poster.Images)
	require.NoError(t, v3().Migrate(db))

	secondImageID := uuid.NewString() + ".jpg"
	require.NoError(t, db.Create(&model.PosterImage{ID: secondImageID, PosterID: poster.ID}).Error)
	require.NoError(t, db.Delete(&poster).Error)
	var count int64
	require.NoError(t, db.Model(&model.PosterImage{}).Count(&count).Error)
	require.Zero(t, count)
}

func TestV3CopiesImages(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Poster{}))
	require.NoError(t, db.Exec("ALTER TABLE posters ADD COLUMN image_id TEXT NOT NULL").Error)
	festival := model.Festival{ID: uuid.New(), Name: "test"}
	require.NoError(t, db.Create(&festival).Error)
	posterID := uuid.New()
	imageID := uuid.NewString() + ".webp"
	require.NoError(t, db.Exec(
		"INSERT INTO posters (id, festival_id, poster_name, description, status, image_id) VALUES (?, ?, ?, ?, ?, ?)",
		posterID, festival.ID, "test", "location", "uncollected", imageID,
	).Error)

	require.NoError(t, v3().Migrate(db))
	var image model.PosterImage
	require.NoError(t, db.First(&image, "id = ?", imageID).Error)
	require.Equal(t, posterID, image.PosterID)
	require.False(t, db.Migrator().HasColumn(&model.Poster{}, "image_id"))
}

func TestMigrateFreshSchemaIncludesPosterImages(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	init, err := Migrate(db)
	require.NoError(t, err)
	require.True(t, init)
	require.True(t, db.Migrator().HasTable(&model.PosterImage{}))
	var ids []string
	require.NoError(t, db.Table("migrations").Order("id").Pluck("id", &ids).Error)
	require.Equal(t, []string{"1", "2", "3", "SCHEMA_INIT"}, ids)
	init, err = Migrate(db)
	require.NoError(t, err)
	require.False(t, init)
}

func TestV3CopyFailureKeepsLegacyImages(t *testing.T) {
	t.Parallel()
	db := migrationTestDB(t)
	require.NoError(t, db.AutoMigrate(&model.Poster{}))
	require.NoError(t, db.Exec("ALTER TABLE posters ADD COLUMN image_id TEXT NOT NULL").Error)
	festival := model.Festival{ID: uuid.New(), Name: "test"}
	require.NoError(t, db.Create(&festival).Error)
	posterIDs := []uuid.UUID{uuid.New(), uuid.New()}
	imageID := uuid.NewString() + ".webp"
	for _, id := range posterIDs {
		require.NoError(t, db.Exec(
			"INSERT INTO posters (id, festival_id, poster_name, description, status, image_id) VALUES (?, ?, ?, ?, ?, ?)",
			id, festival.ID, id.String(), "location", "uncollected", imageID,
		).Error)
	}

	// 同じ画像IDを別のポスターにも紐づけようとすると、主キー制約で失敗する。
	require.Error(t, v3().Migrate(db))
	require.True(t, db.Migrator().HasColumn(&model.Poster{}, "image_id"))
	var imageIDs []string
	require.NoError(t, db.Table("posters").Pluck("image_id", &imageIDs).Error)
	require.Equal(t, []string{imageID, imageID}, imageIDs)

	require.NoError(t, db.Exec("UPDATE posters SET image_id = ? WHERE id = ?", uuid.NewString()+".jpg", posterIDs[1]).Error)
	require.NoError(t, v3().Migrate(db))
	var count int64
	require.NoError(t, db.Model(&model.PosterImage{}).Count(&count).Error)
	require.EqualValues(t, 2, count)
}
