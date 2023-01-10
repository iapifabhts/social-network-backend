package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/iapifabhts/social-network-backend/controller"
	"github.com/iapifabhts/social-network-backend/database"
	_ "github.com/iapifabhts/social-network-backend/docs"
	"github.com/iapifabhts/social-network-backend/repository"
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
	db := database.New(dbConfig)
	userRepository := repository.NewUserRepository(db)
	userController := controller.NewUserController(userRepository)
	app.Post("/signIn", userController.SignIn)
	app.Post("/signUp", userController.SignUp)
	app.Get("/signOut", userController.SignOut)
	app.Get("/users", userController.AllUsers)
	app.Get("/users/:userID", userController.UserByID)
	app.Delete("/users/:userID", userController.DeleteUser)
	app.Patch("/users/:userID", userController.UpdateUser)
}

func main() {
	app := fiber.New()
	routesInit(app)
	log.Fatal(app.Listen(":80"))
}
