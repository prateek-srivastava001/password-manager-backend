package controllers

import (
	"net/http"
	"password-manager/internal/models"

	"github.com/labstack/echo/v4"
)

func Login(ctx echo.Context) error {
	var payload models.LoginRequest

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "works",
		"status":  "success",
	})

}
