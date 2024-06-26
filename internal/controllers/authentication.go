package controllers

import (
	"errors"
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
		Name:        payload.Name,
		Email:       payload.Email,
		Password:    string(hashedPassword),
		Credentials: []models.Credential{},
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

	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "Invalid email or password",
			"status":  "failed",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ctx.JSON(http.StatusConflict, map[string]string{
				"message": "Invalid password",
				"status":  "fail",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
			"status":  "error",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("mostestsecretkeyever"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to generate token",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
		"name":    user.Name,
		"token":   tokenString,
		"status":  "success",
	})
}
