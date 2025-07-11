package database

import "database/sql"

// Migration interface que define los métodos que debe implementar cada migración
type Migration interface {
	Up(db *sql.DB) error
	Down(db *sql.DB) error
	GetName() string
	GetTimestamp() string
}

// BaseMigration estructura base que pueden embeber las migraciones
type BaseMigration struct {
	Name      string
	Timestamp string
}

func (m *BaseMigration) GetName() string {
	return m.Name
}

func (m *BaseMigration) GetTimestamp() string {
	return m.Timestamp
}
