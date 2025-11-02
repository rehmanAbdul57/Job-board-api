package main

import (
	"example.com/job-board/config"
	"example.com/job-board/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	config.ConnectDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
