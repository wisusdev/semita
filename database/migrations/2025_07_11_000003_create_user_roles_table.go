package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreateUserRolesTable struct {
	database.BaseMigration
}

func NewCreateUserRolesTable() *CreateUserRolesTable {
	return &CreateUserRolesTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_user_roles_table",
			Timestamp: "2025_07_11_000003",
		},
	}
}

func (m *CreateUserRolesTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE user_roles (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			role_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			UNIQUE KEY unique_user_role (user_id, role_id),
			INDEX idx_user_roles_user_id (user_id),
			INDEX idx_user_roles_role_id (role_id)
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateUserRolesTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_roles")
	return err
}
