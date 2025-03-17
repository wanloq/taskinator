package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunMigrations executes database migrations using golang-migrate
func RunMigrations() error {
	migrationsPath := "db/migrations"

	cmd := exec.Command("migrate", "-database", os.Getenv("DATABASE_URL"), "-path", migrationsPath, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Migration failed: %s\nError: %v", string(output), err)
		return err
	}

	fmt.Println("✅ Migrations applied successfully")
	return nil
}
