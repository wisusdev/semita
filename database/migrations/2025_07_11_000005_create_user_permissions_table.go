package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreateUserPermissionsTable struct {
	database.BaseMigration
}

func NewCreateUserPermissionsTable() *CreateUserPermissionsTable {
	return &CreateUserPermissionsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_user_permissions_table",
			Timestamp: "2025_07_11_000005",
		},
	}
}

func (m *CreateUserPermissionsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE user_permissions (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			permission_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
			UNIQUE KEY unique_user_permission (user_id, permission_id),
			INDEX idx_user_permissions_user_id (user_id),
			INDEX idx_user_permissions_permission_id (permission_id)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateUserPermissionsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_permissions")
	return err
}
