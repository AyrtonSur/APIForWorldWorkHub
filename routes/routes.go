package routes

import (
	"github.com/gin-gonic/gin"
	"example/APIForWorldWorkHub/controllers"
	"example/APIForWorldWorkHub/middlewares"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)

	router.POST("/users", controllers.Register)

	auth := router.Group("/")
	auth.Use(middlewares.Authenticate())
	{
		auth.GET("/users", controllers.GetUsers)
		auth.GET("/users/:id", controllers.GetUser)
		auth.POST("/services", controllers.AddService)
	}
}