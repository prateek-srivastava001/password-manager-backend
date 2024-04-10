package routes

import (
	"password-manager/internal/controllers"
	"password-manager/internal/middleware"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(incomingRoutes *echo.Echo) {
	incomingRoutes.POST("/signup", controllers.Signup)
	incomingRoutes.POST("/login", controllers.Login)

	auth := incomingRoutes.Group("/api")
	auth.Use(middleware.JWTMiddleware)
	auth.POST("/credential", controllers.AddCredential)
	auth.GET("/all/credentials", controllers.GetAllCredentials)
	auth.PATCH("/credential/:id", controllers.EditCredential)
	auth.DELETE("/credential/:id", controllers.DeleteCredential)
}
