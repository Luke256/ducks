package main

import (
	"os"
	"log/slog"

	repository "github.com/Luke256/ducks/repository/gorm"
	"github.com/Luke256/ducks/router"
	"github.com/Luke256/ducks/router/v1"

	dsnConfig "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")

	router := setup()
	router.Setup()

	if err := router.Start(); err != nil {
		slog.Error("failed to start server:", slog.String("error", err.Error()))
	}
}

func setup() *router.Router {
	dbUser := os.Getenv("NS_MARIADB_USERNAME")
	dbPassword := os.Getenv("NS_MARIADB_PASSWORD")
	dbHost := os.Getenv("NS_MARIADB_HOSTNAME")
	dbPort := os.Getenv("NS_MARIADB_PORT")
	dbName := os.Getenv("NS_MARIADB_DATABASE")

	e := echo.New()

	DSNConfig := dsnConfig.Config{
		User:                 dbUser,
		Passwd:               dbPassword,
		Net:                  "tcp",
		Addr:                 dbHost + ":" + dbPort,
		DBName:               dbName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: DSNConfig.FormatDSN(),
	}))
	if err != nil {
		e.Logger.Fatal("failed to connect database:", err)
	}

	repo := repository.NewGormRepository(db)

	v1Handler := v1.NewHandler(&repo)

	router := router.NewRouter(e, v1Handler, repo)

	return router
}