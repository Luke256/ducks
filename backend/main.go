package main

import (
	"log/slog"
	"os"

	repository "github.com/Luke256/ducks/repository/gorm"
	"github.com/Luke256/ducks/router"
	v1 "github.com/Luke256/ducks/router/v1"
	"github.com/Luke256/ducks/service/festival"
	"github.com/Luke256/ducks/service/poster"
	"github.com/Luke256/ducks/utils/storage/s3"

	dsnConfig "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	bucketName := os.Getenv("S3_BUCKET_NAME")

	e := echo.New()

	// address CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

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
		slog.Error("failed to connect database:", slog.String("error", err.Error()))
		panic(err)
	}

	repo, _, err := repository.NewGormRepository(db, true)
	if err != nil {
		slog.Error("failed to create repository:", slog.String("error", err.Error()))
		panic(err)
	}

	storage, err := s3.NewS3Storage(bucketName)
	if err != nil {
		slog.Error("failed to create s3 storage:", slog.String("error", err.Error()))
		panic(err)
	}

	festivalManager := festival.NewManagerImpl(repo)
	posterManager := poster.NewManagerImpl(repo, storage)

	v1Handler := v1.NewHandler(repo, festivalManager, posterManager, storage)

	router := router.NewRouter(e, v1Handler, repo)

	return router
}
