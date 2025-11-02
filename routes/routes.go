package routes

import (
	"example.com/job-board/controllers"
	"example.com/job-board/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/signup", controllers.Signup)
	api.Post("/login", controllers.Login)
	api.Get("/jobs", controllers.GetJobs)

	authenticated := api.Group("")
	authenticated.Use(middlewares.JWTAuthMiddleware())
	authenticated.Post("/jobs", controllers.CreateJob)
	authenticated.Post("/apply", controllers.ApplyToJob)
	authenticated.Get("/job/:id/applications", controllers.GetApplicationsForJob)
	authenticated.Get("/employer/jobs", controllers.GetMyJobs)

}
