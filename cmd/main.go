package main

import (
	"context"
	"embed"
	"io/fs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mwozniak/missing-aed/internal/client"
	"github.com/mwozniak/missing-aed/internal/controller"
	"github.com/mwozniak/missing-aed/internal/dto"
	"github.com/mwozniak/missing-aed/internal/model"
	"github.com/mwozniak/missing-aed/internal/repository"
	"github.com/mwozniak/missing-aed/internal/service"
)

//go:embed app
var appFiles embed.FS

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	if err := godotenv.Load(); err != nil {
		logrus.Info("no .env file loaded")
	}

	config := dto.NewConfig()

	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		logrus.Panic(err)
	}
	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		logrus.Panic(err)
	}

	clients := client.NewClients(config)
	repositories := repository.NewRepositories(db)
	services, err := service.NewServices(clients, repositories, config)
	if err != nil {
		logrus.Panic(err)
	}
	controllers := controller.NewControllers(services)

	ctx := context.Background()
	go services.Osm().Start(ctx)

	e := echo.New()
	controllers.Route(e)

	staticFs, err := fs.Sub(appFiles, "app")
	if err != nil {
		logrus.Panic(err)
	}
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Index:      "index.html",
		HTML5:      true,
		Filesystem: staticFs,
	}))

	logrus.Info("starting server on port 3000")
	logrus.Fatal(e.Start(":3000"))
}
