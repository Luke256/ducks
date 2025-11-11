package v1

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Luke256/ducks/migration"
	"github.com/Luke256/ducks/repository"
	gormRepo "github.com/Luke256/ducks/repository/gorm"
	"github.com/Luke256/ducks/service/festival"
	"github.com/Luke256/ducks/service/poster"
	"github.com/Luke256/ducks/utils"
	mockstorage "github.com/Luke256/ducks/utils/storage/mock_storage"
	"github.com/labstack/echo/v4"

	driverMysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbPrefix = "traq-ducks-router-test-"
	common   = "common"
	s1       = "s1"
)

var (
	envs = map[string]*env{}
)

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
		env := &env{}
		dbConfig := *config
		dbConfig.DBName = fmt.Sprintf("%s%s", dbPrefix, key)

		// DB接続
		engine, err := gorm.Open(mysql.New(mysql.Config{
			DSN: dbConfig.FormatDSN(),
		}))
		if err != nil {
			panic(err)
		}
		env.DB = engine
		db, err := engine.DB()
		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(20)
		if err := migration.DropAll(engine); err != nil {
			panic(err)
		}

		repo, _, err := gormRepo.NewGormRepository(engine, true)
		if err != nil {
			panic(err)
		}
		env.Repo = repo

		env.Storage = &mockstorage.MockStorage{}

		env.FM = festival.NewManagerImpl(repo)
		env.PM = poster.NewManagerImpl(repo, env.Storage)

		// サーバー
		e := echo.New()
		e.HideBanner = true
		e.HidePort = true

		handlers := NewHandler(
			repo,
			env.FM,
			env.PM,
		)
		handlers.Setup(e.Group("/api"))
		env.Server = httptest.NewServer(e)

		envs[key] = env
	}

	code := m.Run()

	for _, e := range envs {
		db, err := e.DB.DB()
		if err != nil {
			panic(err)
		}
		e.Server.Close()
		db.Close()
	}

	os.Exit(code)
}

type env struct {
	Server  *httptest.Server
	DB      *gorm.DB
	Repo    repository.Repository
	FM      festival.Manager
	PM      poster.Manager
	Storage *mockstorage.MockStorage
}
