package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreatePermissionsTable struct {
	database.BaseMigration
}

func NewCreatePermissionsTable() *CreatePermissionsTable {
	return &CreatePermissionsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_permissions_table",
			Timestamp: "2025_07_11_000002",
		},
	}
}

func (m *CreatePermissionsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE permissions (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) UNIQUE NOT NULL,
			guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
			description TEXT DEFAULT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_permissions_name (name),
			INDEX idx_permissions_guard_name (guard_name)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreatePermissionsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS permissions")
	return err
}
