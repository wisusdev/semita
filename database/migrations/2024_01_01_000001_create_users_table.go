package migrations

import (
	"database/sql"
	"web_utilidades/database"
)

type CreateUsersTable struct {
	database.BaseMigration
}

func NewCreateUsersTable() *CreateUsersTable {
	return &CreateUsersTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_users_table",
			Timestamp: "2024_01_01_000001",
		},
	}
}

func (m *CreateUsersTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			email_verified_at DATETIME DEFAULT NULL,
			remember_token VARCHAR(100) DEFAULT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateUsersTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}
