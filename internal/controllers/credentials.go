package controllers

import (
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/models"

	"github.com/labstack/echo/v4"
)

func AddCredential(ctx echo.Context) error {
	var credential models.Credential
	
	if err := ctx.Bind(&credential); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "failed",
		})
	}

	email := "prateek@gmail.com"
	err := database.AddCredential(email, credential)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to add credential",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "Credential added successfully",
		"status":  "success",
	})
}
