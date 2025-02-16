package routes

import (
	"example/APIForWorldWorkHub/controllers"
	"example/APIForWorldWorkHub/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)
	router.POST("/refresh-token", controllers.RefreshToken)
	router.POST("/email-exists", controllers.CheckEmailExists)
	router.POST("/users", controllers.Register)
	router.POST("/zip-exists", controllers.CheckZip)

	auth := router.Group("/")
	auth.Use(middlewares.Authenticate())
	{
		auth.POST("/logout", controllers.Logout)
		auth.GET("/users", middlewares.Authorize("view_users"), controllers.GetUsers)
		auth.GET("/users/:id", middlewares.Authorize("view_user"), controllers.GetUser)
		auth.POST("/services", middlewares.Authorize("create_service"), controllers.AddService)
		auth.PATCH("/users/:id", middlewares.AuthorizeSelfOrPermission("update_user"), controllers.UpdateUser)
		auth.DELETE("/users/:id", middlewares.Authorize("delete_user"), controllers.DeleteUser)
		auth.GET("/services", middlewares.Authorize("view_services"), controllers.GetServices)
		auth.POST("/users/create", middlewares.Authorize("create_user"), controllers.CreateUser)
		auth.POST("/roles", middlewares.Authorize("create_role"), controllers.CreateRole)
		auth.GET("/roles", middlewares.Authorize("view_roles"), controllers.GetRoles)
		auth.DELETE("/roles", middlewares.Authorize("delete_role"), controllers.DeleteRole)
		auth.GET("/current-user", controllers.GetCurrentUser)
	}
}
