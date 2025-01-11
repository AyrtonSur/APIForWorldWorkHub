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
		auth.GET("/users", middlewares.Authorize("view_users"), controllers.GetUsers)
		auth.GET("/users/:id", middlewares.Authorize("view_user"), controllers.GetUser)
		auth.POST("/services", middlewares.Authorize("create_service"), controllers.AddService)
		auth.PATCH("/users/:id", middlewares.AuthorizeSelfOrPermission("update_user"), controllers.UpdateUser)
		auth.DELETE("/users/:id", middlewares.Authorize("delete_user"), controllers.DeleteUser)
		auth.GET("/services", middlewares.Authorize("view_services"), controllers.GetServices)
	}
}