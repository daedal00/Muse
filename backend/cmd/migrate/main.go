package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daedal00/muse/backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|drop|version|force <version>]")
	}

	command := os.Args[1]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		cfg.GetDatabaseDSN(),
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("✅ Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("✅ Migrations rolled back successfully")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("Failed to drop database: %v", err)
		}
		fmt.Println("✅ Database dropped successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d, dirty: %t\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go force <version>")
		}
		version := os.Args[2]
		var v int
		if _, err := fmt.Sscanf(version, "%d", &v); err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		fmt.Printf("✅ Forced migration to version %d\n", v)

	default:
		log.Fatal("Unknown command. Use: up, down, drop, version, or force")
	}
}
