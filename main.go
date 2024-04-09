package main

import (
	"log"
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Password Manager",
			"status":  "true",
		})
	})

	routes.AuthRoutes(app)

	mongoURI := "mongodb+srv://prateeksrivastava201:ozr6FDh2RUMO5qbd@cluster0.kg1bllf.mongodb.net/?retryWrites=true&w=majority"

	if err := database.Connect(mongoURI); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	app.Logger.Fatal(app.Start(":8080"))
}
