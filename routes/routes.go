package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"example/APIForWorldWorkHub/controllers"
	"example/APIForWorldWorkHub/middlewares"
)

func SetupRoutes(router *gin.Engine, validate *validator.Validate) {
	router.POST("/login", func(c *gin.Context) {
		controllers.Login(c, validate)
	})

	router.POST("/users", func(c *gin.Context) {
		controllers.Register(c, validate)
	})

	auth := router.Group("/")
	auth.Use(middlewares.Authenticate())
	{
		auth.GET("/users", controllers.GetUsers)

		auth.GET("/users/:id", controllers.GetUser)
		auth.POST("/services", func(c *gin.Context) {
			controllers.AddService(c, validate)
		})
	}
}