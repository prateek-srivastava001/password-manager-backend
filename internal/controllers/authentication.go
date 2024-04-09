package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prateek-srivastava001/password-manager-backend/internal/models"
)

func Login(ctx echo.Context) error {
	var payload models.LoginRequest

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": error.Error(),
			"status":  "failed",
		})
	}

	
}
