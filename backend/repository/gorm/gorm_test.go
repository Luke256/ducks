package gorm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Luke256/ducks/migration"
	"github.com/Luke256/ducks/model"
	"github.com/Luke256/ducks/utils"
	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbPrefix = "traq-ducks-test-"
	common = "common"
	s1 = "s1"
)

var (
	repositories = map[string]*GormRepository{}
)

// 前処理
func TestMain(m *testing.M) {
	dbUser := utils.GetEnvOrDefault("NS_MARIADB_USERNAME", "root")
	dbPassword := utils.GetEnvOrDefault("NS_MARIADB_PASSWORD", "password")
	dbHost := utils.GetEnvOrDefault("NS_MARIADB_HOST", "localhost")
	dbPort := utils.GetEnvOrDefault("NS_MARIADB_PORT", "3307")
	dbs := []string{
		common, s1,
	}
	config := &driverMysql.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	if err := migration.CreateDatabasesIfNotExists("mysql", config.FormatDSN(), dbPrefix, dbs...); err != nil {
		panic(err)
	}

	for _, key := range dbs {
		dbConfig := *config
		dbConfig.DBName = fmt.Sprintf("%s%s", dbPrefix, key)
		
		engine, err := gorm.Open(mysql.New(mysql.Config{
			DSN: dbConfig.FormatDSN(),
		}))
		if err != nil {
			panic(err)
		}

		db, err := engine.DB()
		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(20)

		if err := migration.DropAll(engine); err != nil {
			panic(err)
		}

		repo, _, err := NewGormRepository(engine, true)
		if err != nil {
			panic(err)
		}

		repositories[key] = repo
	}

	code := m.Run()

	for _, repo := range repositories {
		db, err := repo.db.DB()
		if err != nil {
			panic(err)
		}

		db.Close()
	}
	os.Exit(code)
}

func setup(t *testing.T, dbKey string) *GormRepository {
	repo, ok := repositories[dbKey]
	if !ok {
		t.Fatalf("repository %s not found", dbKey)
	}
	return repo
}

func mustCreateFestival(t *testing.T, repo *GormRepository, name string, description string) model.Festival {
	t.Helper()

	festival, err := repo.RegisterFestival(name, description)
	if err != nil {
		t.Fatalf("failed to register festival: %v", err)
	}
	return festival
}

func mustCreatePoster(t *testing.T, repo *GormRepository, festivalID uuid.UUID, posterName string, description string, imageID string) model.Poster {
	t.Helper()

	poster, err := repo.RegisterPoster(festivalID, posterName, description, imageID)
	if err != nil {
		t.Fatalf("failed to register poster: %v", err)
	}

	return poster
}

func mustCreateStockItem(t *testing.T, repo *GormRepository, name string, description string, category string, imageID string) model.StockItem {
	t.Helper()

	item, err := repo.RegisterStockItem(name, description, category, imageID)
	if err != nil {
		t.Fatalf("failed to register stock item: %v", err)
	}

	return item
}