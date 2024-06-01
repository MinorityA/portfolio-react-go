package main

import (
	"mongorest/common"
	"mongorest/router"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {

	err := common.LoadEnv()
	if err != nil {
		return err
	}

	err = common.InitDB()
	if err != nil {
		return err
	}

	defer common.CloseDB()

	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddBlogGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "4000"
	}
	app.Listen(":" + port)

	return nil
}
