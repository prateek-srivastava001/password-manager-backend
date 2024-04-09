package controllers

import (
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx echo.Context) error {
	var payload models.SignUpRequest

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
			"status":  "failed",
		})
	}

	if _, err := database.GetUserByEmail(payload.Email); err == nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"message": "Email already exists",
			"status":  "failed",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to hash password",
			"status":  "failed",
		})
	}

	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}

	if err := database.CreateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create user",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"message": "User created successfully",
		"status":  "success",
	})
}

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
