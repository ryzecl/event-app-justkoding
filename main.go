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
	config.InitImageKit()

	server := gin.Default()

	// Route
	api:= server.Group("/api")
	{
		
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventbyId)
		

		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.RequiredAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)

			protected.POST("/events", controllers.CreateEvent)
			protected.PUT("/events/:id", controllers.UpdateEvent)
			protected.DELETE("/events/:id", controllers.DeleteEvent)
		}
	}

	server.Run(":8080")
}