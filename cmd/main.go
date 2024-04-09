package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message":"Welcome to Password Manager",
			"status": "true",
		})
	})
	app.Logger.Fatal(app.Start(":8080"))
}
