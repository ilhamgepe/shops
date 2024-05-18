package server

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamgepe/shops/packages/config"
	"github.com/ilhamgepe/shops/packages/logs"
	"github.com/jmoiron/sqlx"
)

type FiberServer struct {
	*fiber.App
	*validator.Validate
	DB *sqlx.DB
}

func New(db *sqlx.DB) *FiberServer {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ServerHeader: fmt.Sprintf("x-%s", config.Get.APP_NAME),
		AppName:      config.Get.APP_NAME,
	})
	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &logs.Logger}))

	return &FiberServer{
		App:      app,
		Validate: validator.New(),
		DB:       db,
	}
}

func (fs *FiberServer) Run() error {
	return fs.App.Listen(fmt.Sprintf(":%s", config.Get.APP_PORT))
}

func (fs *FiberServer) Shutdown() error {
	return fs.App.Shutdown()
}
