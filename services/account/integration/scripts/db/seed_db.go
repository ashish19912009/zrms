package db

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SeedDB() error {
	// seed logic
	fmt.Println("Seeding database...")
	return nil
}

func main() {
	// Parse environment flag
	env := flag.String("env", "local", "Environment to use: local, staging, prod")
	if *env == "prod" {
		log.Fatal("Never seed production DB from local script")
	}
	flag.Parse()

	// Load .env file based on the environment
	envFile := fmt.Sprintf(".env.%s", *env)
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("failed to load %s: %v", envFile, err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL not set in %s", envFile)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(`TRUNCATE TABLE users.accounts RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatalf("failed to truncate accounts: %v", err)
	}

	now := time.Now().UTC()
	_, err = tx.Exec(`
		INSERT INTO users.accounts (id, mobile_no, name, role, status, employee_id, created_at)
		VALUES
			(gen_random_uuid(), '7000000001', 'Admin User', 'admin', 'active', 'EMP1001', $1),
			(gen_random_uuid(), '7000000002', 'Delivery Guy', 'delivery', 'active', 'EMP1002', $1),
			(gen_random_uuid(), '7000000003', 'Support Agent', 'support', 'inactive', 'EMP1003', $1)
	`, now)
	if err != nil {
		log.Fatalf("failed to insert seed data: %v", err)
	}

	fmt.Printf("✅ Seeded %s environment database successfully\n", *env)
}
