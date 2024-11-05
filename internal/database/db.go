package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"password-manager/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

var DB *sql.DB

func ConnectPostgres() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to NeonDB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Cannot ping NeonDB: %v", err)
	}

	fmt.Println("Successfully connected to NeonDB")
}

func DisconnectPostgres() {
	if DB != nil {
		DB.Close()
	}
}

func RunMigrations() {
	if DB == nil {
		log.Fatal("Database connection is not initialized")
	}

	migrationsDir := "./db/migrations"
	if err := goose.Up(DB, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully.")
}

func GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1"
	row := DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %v", err)
	}

	return &user, nil
}

func CreateUser(user models.User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := DB.Exec(query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func AddCredential(email string, credential models.Credential) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO credentials (id, user_id, name, username, password, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = DB.Exec(query, credential.ID, user.ID, credential.Name, credential.Username, credential.Password, credential.URL, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error adding credential: %v", err)
	}

	return nil
}

func GetCredentialsByEmail(email string) ([]models.Credential, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %v", err)
	}

	query := "SELECT id, user_id, name, username, password, url, created_at, updated_at FROM credentials WHERE user_id = $1"
	rows, err := DB.Query(query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching credentials: %v", err)
	}
	defer rows.Close()

	var credentials []models.Credential
	for rows.Next() {
		var credential models.Credential
		if err := rows.Scan(&credential.ID, &credential.UserID, &credential.Name, &credential.Username, &credential.Password, &credential.URL, &credential.CreatedAt, &credential.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning credential: %v", err)
		}
		credentials = append(credentials, credential)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(credentials) == 0 {
		return nil, fmt.Errorf("credentials not found")
	}

	return credentials, nil
}

func UpdateCredential(email string, credential models.Credential) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("error fetching user by email: %v", err)
	}

	query := `
		UPDATE credentials
		SET name = $1, username = $2, password = $3, url = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`
	_, err = DB.Exec(query, credential.Name, credential.Username, credential.Password, credential.URL, time.Now(), credential.ID, user.ID)
	if err != nil {
		return fmt.Errorf("error updating credential: %v", err)
	}

	return nil
}

func DeleteCredential(email string, credentialID string) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("error fetching user by email: %v", err)
	}

	query := `
		DELETE FROM credentials
		WHERE id = $1 AND user_id = $2
	`
	_, err = DB.Exec(query, credentialID, user.ID)
	if err != nil {
		return fmt.Errorf("error deleting credential: %v", err)
	}

	return nil
}
