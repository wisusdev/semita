package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreateRolesTable struct {
	database.BaseMigration
}

func NewCreateRolesTable() *CreateRolesTable {
	return &CreateRolesTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_roles_table",
			Timestamp: "2025_07_11_000001",
		},
	}
}

func (m *CreateRolesTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE roles (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) UNIQUE NOT NULL,
			guard_name VARCHAR(255) NOT NULL DEFAULT 'web',
			description TEXT DEFAULT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_roles_name (name),
			INDEX idx_roles_guard_name (guard_name)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateRolesTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS roles")
	return err
}
