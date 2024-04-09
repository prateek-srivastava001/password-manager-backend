package routes

import (
	"password-manager/internal/controllers"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(incomingRoutes *echo.Echo) {
	incomingRoutes.POST("/signup", controllers.Signup)
}
