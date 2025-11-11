package gorm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Luke256/ducks/migration"
	driverMysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbPrefix = "ducks-test-"
	common = "common"
)

var (
	repositories = map[string]*GormRepository{}
)

// 前処理
func TestMain(m *testing.M) {
	dbUser := os.Getenv("NS_MARIADB_USERNAME")
	dbPassword := os.Getenv("NS_MARIADB_PASSWORD")
	dbHost := os.Getenv("NS_MARIADB_HOST")
	dbPort := os.Getenv("NS_MARIADB_PORT")
	dbs := []string{
		common,
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

	m.Run()

	for _, repo := range repositories {
		db, err := repo.db.DB()
		if err != nil {
			panic(err)
		}

		db.Close()
	}
}

func setup(t *testing.T, dbKey string) *GormRepository {
	repo, ok := repositories[dbKey]
	if !ok {
		t.Fatalf("repository %s not found", dbKey)
	}
	return repo
}