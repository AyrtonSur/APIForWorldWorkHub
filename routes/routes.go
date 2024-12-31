package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"example/APIForWorldWorkHub/controllers"
)

func SetupRoutes(router *gin.Engine, validate *validator.Validate) {
	router.GET("/users", controllers.GetUsers)
	router.POST("/users", func(c *gin.Context) {
		controllers.Register(c, validate)
	})
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/services", func(c *gin.Context) {
		controllers.AddService(c, validate)
	})
}