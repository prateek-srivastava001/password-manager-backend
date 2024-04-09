package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/prateek-srivastava001/password-manager-backend/internal/models"
)

func Login(ctx echo.Context) error {
	var payload models.LoginRequest
}
