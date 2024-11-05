	package models

	import (
		"time"
	)

	type User struct {
		ID        string    `json:"id" db:"id"`
		Name      string    `json:"name" db:"name"`
		Email     string    `json:"email" db:"email"`
		Password  string    `json:"password" db:"password"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	type Credential struct {
		ID        string    `json:"id" db:"id"`
		UserID    string    `json:"user_id" db:"user_id"`
		Name      string    `json:"name" db:"name"`
		Username  string    `json:"username" db:"username"`
		Password  string    `json:"password" db:"password"`
		URL       string    `json:"url" db:"url"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	type SignUpRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
