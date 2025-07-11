package migrations

import (
	"database/sql"
	"web_utilidades/database"
)

type CreateOAuthScopesTable struct {
	database.BaseMigration
}

func NewCreateOAuthScopesTable() *CreateOAuthScopesTable {
	return &CreateOAuthScopesTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_oauth_scopes_table",
			Timestamp: "2025_07_06_000003",
		},
	}
}

func (m *CreateOAuthScopesTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE oauth_scopes (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL UNIQUE,
			description VARCHAR(255),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateOAuthScopesTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_scopes")
	return err
}
