package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("--User Service Migrate Start--")
	dbURL := "postgres://user:password@localhost:5442/mydatabase?sslmode=disable"
	m, err := migrate.New(
		"file://files",
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed create migrate instance : %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully")

}
