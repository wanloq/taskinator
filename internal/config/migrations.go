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
	migratePath := "migrate"
	// migratePath := "/usr/local/bin/migrate"

	cmd := exec.Command(migratePath, "-database", os.Getenv("DATABASE_URL"), "-path", migrationsPath, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("DB Migration failed: %s\nError: %v", string(output), err)
		return err
	}

	fmt.Println("✅ Migrations applied successfully")
	return nil
}
