package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prateek-srivastava001/password-manager-backend/internal/database"
)

func main() {
	app := echo.New()
	app.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Password Manager",
			"status":  "true",
		})
	})

	mongoURI := "mongodb+srv://prateeksrivastava201:hrXbQquD6WJgvU6w@cluster0.qxkjh9k.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	if err := database.Connect(mongoURI); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	app.Logger.Fatal(app.Start(":8080"))
}
