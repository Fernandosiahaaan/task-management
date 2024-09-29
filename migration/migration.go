package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("--Microservice Migrate Start--")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbURL := os.Getenv("POSTGRES_URI")
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
