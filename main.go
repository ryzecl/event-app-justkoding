package main

import (
	"log"

	"example.com/event-app/config"
	"example.com/event-app/controllers"
	"example.com/event-app/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()

	// Route
	api:= server.Group("/api")
	{
		api.POST("/events", controllers.CreateEvent)
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventbyId)
		api.PUT("/events/:id", controllers.UpdateEvent)
		api.DELETE("/events/:id", controllers.DeleteEvent)

		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.RequiredAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)
		}
	}

	server.Run(":8080")
}