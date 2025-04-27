package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func RunMigrations(databaseURL string) {
	dbInstance, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("failed to acquire database connection:", err)
	}
	defer dbInstance.Close()

	driver, err := postgres.WithInstance(dbInstance, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to create migrate driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("failed to create migrate instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed:", err)
	}

	fmt.Println("Database migrated successfully")
}
