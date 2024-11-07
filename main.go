package main

import (
	"log"
	"net/http"
	"os"
	"password-manager/internal/database"
	"password-manager/internal/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.ConnectPostgres()
	defer database.DisconnectPostgres()
	database.RunMigrations()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	routes.AuthRoutes(e)

	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Password Manager",
			"status":  "true",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}
