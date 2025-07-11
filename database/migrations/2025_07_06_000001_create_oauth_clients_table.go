package migrations

import (
	"database/sql"
	"web_utilidades/database"
)

type CreateOAuthClientsTable struct {
	database.BaseMigration
}

func NewCreateOAuthClientsTable() *CreateOAuthClientsTable {
	return &CreateOAuthClientsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_oauth_clients_table",
			Timestamp: "2025_07_06_000001",
		},
	}
}

func (m *CreateOAuthClientsTable) Up(db *sql.DB) error {
	query := `
        CREATE TABLE oauth_clients (
            id INT PRIMARY KEY AUTO_INCREMENT,
            name VARCHAR(255) NOT NULL,
            client_id VARCHAR(100) NOT NULL UNIQUE,
            client_secret VARCHAR(255) NOT NULL,
            redirect_uri VARCHAR(255),
            grant_types VARCHAR(255),
            scopes VARCHAR(255),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        )
    `
	_, err := db.Exec(query)
	return err
}

func (m *CreateOAuthClientsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_clients")
	return err
}
