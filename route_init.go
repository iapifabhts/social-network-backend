package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/iapifabhts/social-network-backend/controller"
	"github.com/iapifabhts/social-network-backend/database"
	_ "github.com/iapifabhts/social-network-backend/docs"
	"github.com/iapifabhts/social-network-backend/middleware"
	"github.com/iapifabhts/social-network-backend/repository"
	"github.com/iapifabhts/social-network-backend/session"
)

func routeInit(app *fiber.App) {
	db := database.New()
	sessionStore := session.New()
	sessionChecker := middleware.NewSessionChecker(sessionStore)
	userRepo := repository.NewUserRepository(db)
	userController := controller.NewUserController(userRepo, sessionStore)
	app.Get("/swagger-ui/*", swagger.New(swagger.Config{
		CustomStyle: "* {margin: 0;padding: 0;box-sizing: border-box;}",
	}))
	app.Post("/signIn", userController.SignIn)
	app.Post("/signUp", userController.SignUp)
	app.Get("/signOut", sessionChecker.Check, userController.SignOut)
	app.Get("/users", sessionChecker.Check, userController.AllUsers)
	app.Get("/users/:userID", sessionChecker.Check, userController.UserByID)
	app.Patch("/users/:userID", sessionChecker.Check, userController.UpdateUser)
	app.Delete("/users/:userID", sessionChecker.Check, userController.DeleteUser)
}
