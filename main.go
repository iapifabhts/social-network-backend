package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var dbConfig string

func init() {
	godotenv.Load()
	dbConfig = fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"host=%s "+
			"port=%s "+
			"dbname=%s "+
			"sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func routesInit(app *fiber.App) {
	app.Get("/swagger-ui/*", swagger.New(swagger.Config{
		CustomStyle: "* {margin: 0;padding: 0;box-sizing: border-box;}",
	}))
}

func main() {
	app := fiber.New()
	routesInit(app)
	log.Fatal(app.Listen(":80"))
}
