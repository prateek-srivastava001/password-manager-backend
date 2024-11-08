package controllers

import (
	"fmt"
	"net/http"
	"password-manager/internal/database"
	"password-manager/internal/middleware"
	"password-manager/internal/models"
	"time"

	"github.com/google/uuid"
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
		ID:        uuid.NewString(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
			"message": "Invalid login details",
			"status":  "failed",
		})
	}

	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid credentials",
			"status":  "failed",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Invalid credentials",
			"status":  "failed",
		})
	}

	fmt.Println(user.Email)

	token, err := middleware.GenerateJWT(user.Email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to generate token",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"name":    user.Name,
		"token":   token,
	})
}

func AddCredential(ctx echo.Context) error {
	var credential models.Credential

	if err := ctx.Bind(&credential); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "failed",
		})
	}

	credential.ID = uuid.NewString()

	email, ok := ctx.Get("user").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
			"status":  "failed",
		})
	}

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
	email, ok := ctx.Get("user").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
			"status":  "failed",
		})
	}

	creds, err := database.GetCredentialsByEmail(email)
	fmt.Println(err)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to get credentials",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Credentials fetched successfully",
		"status":  "success",
		"data":    creds,
	})
}

func UpdateCredential(ctx echo.Context) error {
	email, ok := ctx.Get("user").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
			"status":  "failed",
		})
	}

	credentialID := ctx.Param("id")

	var updatedCredential models.Credential
	if err := ctx.Bind(&updatedCredential); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
			"status":  "failed",
		})
	}

	updatedCredential.ID = credentialID

	err := database.UpdateCredential(email, updatedCredential)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update credential",
			"status":  "failed",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Credential updated successfully",
		"status":  "success",
	})
}

func DeleteCredential(ctx echo.Context) error {
	email, ok := ctx.Get("user").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": "User not authenticated",
			"status":  "failed",
		})
	}

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
