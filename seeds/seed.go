package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq" // postgres driver
	"zaid-paper-disbursement/config"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Truncate all tables
	if err := truncateTables(db); err != nil {
		log.Fatalf("Error truncating tables: %v", err)
	}

	// Seed users
	userID, err := seedUser(db)
	if err != nil {
		log.Fatalf("Error seeding user: %v", err)
	}

	// Seed user balance
	if err := seedUserBalance(db, userID); err != nil {
		log.Fatalf("Error seeding user balance: %v", err)
	}

	log.Println("Seeding completed successfully")
}

func truncateTables(db *sql.DB) error {
	_, err := db.Exec(`
		TRUNCATE TABLE transactions CASCADE;
		TRUNCATE TABLE user_balances CASCADE;
		TRUNCATE TABLE users CASCADE;
	`)
	return err
}

func seedUser(db *sql.DB) (int64, error) {
	var userID int64
	err := db.QueryRow(`
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`, "Zaid", "zaid@example.com").Scan(&userID)
	return userID, err
}

func seedUserBalance(db *sql.DB, userID int64) error {
	_, err := db.Exec(`
		INSERT INTO user_balances (user_id, balance)
		VALUES ($1, $2)
	`, userID, 1000.00)
	return err
}