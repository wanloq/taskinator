package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	databaseURL := "postgres://admin:adminpassword@localhost:5432/taskinator"

	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer DB.Close()

	fmt.Println("✅ Connected to PostgreSQL")

}

func DBMigrate() {
	// SQL to create tables
	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role VARCHAR(20) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(20) DEFAULT 'pending',
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	// Execute migration
	_, err := DB.Exec(context.Background(), createTables)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}

	fmt.Println("Database migration completed successfully!")
}
