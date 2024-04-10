package controllers

import (
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/models"

	"github.com/google/uuid"
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

	credential.ID = uuid.New().String()

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

func GetAllCredentials(ctx echo.Context) error {
	email := "prateek@gmail.com"
	creds, err := database.GetCredentialsForUser(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch credentials",
			"status":  "failed",
		})
	}
	return ctx.JSON(http.StatusOK, creds)
}

func EditCredential(ctx echo.Context) error {
	email := "prateek@gmail.com"
	credentialID := ctx.Param("id")

	var updatedCredential models.Credential
	if err := ctx.Bind(&updatedCredential); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "failed",
		})
	}

	err := database.EditCredential(email, credentialID, updatedCredential)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to edit credential",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Credential edited successfully",
		"status":  "success",
	})
}

func DeleteCredential(ctx echo.Context) error {
	email := "prateek@gmail.com"
	credentialID := ctx.Param("id")

	err := database.DeleteCredential(email, credentialID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete credential",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Credential deleted successfully",
		"status":  "success",
	})
}
