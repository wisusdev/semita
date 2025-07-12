package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreateOAuthTokensTable struct {
	database.BaseMigration
}

func NewCreateOAuthTokensTable() *CreateOAuthTokensTable {
	return &CreateOAuthTokensTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_oauth_tokens_table",
			Timestamp: "2025_07_06_000002",
		},
	}
}

func (m *CreateOAuthTokensTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE oauth_tokens (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT,
			client_id INT NOT NULL,
			access_token VARCHAR(512) NOT NULL UNIQUE,
			refresh_token VARCHAR(512) NOT NULL UNIQUE,
			scopes VARCHAR(255),
			revoked TINYINT(1) DEFAULT 0,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (client_id) REFERENCES oauth_clients(id)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateOAuthTokensTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_tokens")
	return err
}
