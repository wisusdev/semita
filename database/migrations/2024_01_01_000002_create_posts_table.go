package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type CreatePostsTable struct {
	database.BaseMigration
}

func NewCreatePostsTable() *CreatePostsTable {
	return &CreatePostsTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_posts_table",
			Timestamp: "2024_01_01_000002",
		},
	}
}

func (m *CreatePostsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE posts (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			published BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreatePostsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS posts")
	return err
}
