package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func RunMigrations(databaseURL string) {
	dbInstance, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("failed to acquire database connection:", err)
	}
	defer dbInstance.Close()

	// SQL queries for creating the tables
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    email VARCHAR(255) UNIQUE NOT NULL,
	    created_at TIMESTAMP DEFAULT current_timestamp,
	    updated_at TIMESTAMP DEFAULT current_timestamp
	);`

	createPollsTableSQL := `
	CREATE TABLE IF NOT EXISTS polls (
	    id SERIAL PRIMARY KEY,
	    title TEXT NOT NULL,
	    options TEXT[] NOT NULL,
	    tags TEXT[] NOT NULL,
	    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	    created_at TIMESTAMPTZ DEFAULT now(),
	    updated_at TIMESTAMPTZ DEFAULT now()
	);`

	createUserVotesTableSQL := `
	CREATE TABLE IF NOT EXISTS user_votes (
	    id SERIAL PRIMARY KEY,
	    user_id INT NOT NULL,
	    poll_id INT NOT NULL,
	    option_index INT,
	    skipped BOOLEAN DEFAULT FALSE,
	    created_at TIMESTAMPTZ DEFAULT now(),
	    UNIQUE (user_id, poll_id)
	);`

	// Execute the SQL queries
	_, err = dbInstance.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal("failed to execute SQL for users table:", err)
	}

	_, err = dbInstance.Exec(createPollsTableSQL)
	if err != nil {
		log.Fatal("failed to execute SQL for polls table:", err)
	}

	_, err = dbInstance.Exec(createUserVotesTableSQL)
	if err != nil {
		log.Fatal("failed to execute SQL for user_votes table:", err)
	}

	fmt.Println("Database migrated successfully")
}
