package migrations

import (
	"database/sql"
	"semita/app/core/database"
)

type CreatePasswordResetsTable struct {
	database.BaseMigration
}

func NewCreatePasswordResetsTable() *CreatePasswordResetsTable {
	return &CreatePasswordResetsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_password_resets_table",
			Timestamp: "2025_07_07_000001",
		},
	}
}

func (m *CreatePasswordResetsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE password_resets (
			email VARCHAR(255) NOT NULL,
			token VARCHAR(255) NOT NULL,
			created_at DATETIME NOT NULL,
			PRIMARY KEY (email, token)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreatePasswordResetsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS password_resets")
	return err
}
