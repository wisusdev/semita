package migrations

import (
	"database/sql"
	"log"
)

func CreateSeedersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS seeders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		rollback_at TIMESTAMP NULL,
		INDEX idx_seeder_name (name)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating seeders table: %v", err)
		return err
	}

	log.Println("Seeders table created successfully")
	return nil
}

func DropSeedersTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS seeders;`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error dropping seeders table: %v", err)
		return err
	}

	log.Println("Seeders table dropped successfully")
	return nil
}
