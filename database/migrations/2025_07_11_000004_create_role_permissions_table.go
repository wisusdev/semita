package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreateRolePermissionsTable struct {
	database.BaseMigration
}

func NewCreateRolePermissionsTable() *CreateRolePermissionsTable {
	return &CreateRolePermissionsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_role_permissions_table",
			Timestamp: "2025_07_11_000004",
		},
	}
}

func (m *CreateRolePermissionsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE role_permissions (
			id INT PRIMARY KEY AUTO_INCREMENT,
			role_id INT NOT NULL,
			permission_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
			UNIQUE KEY unique_role_permission (role_id, permission_id),
			INDEX idx_role_permissions_role_id (role_id),
			INDEX idx_role_permissions_permission_id (permission_id)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateRolePermissionsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS role_permissions")
	return err
}
