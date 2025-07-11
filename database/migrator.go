package database

import (
	"database/sql"
	"fmt"
	"sort"
)

type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register registra una nueva migración
func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// CreateMigrationsTable crea la tabla de migraciones si no existe
func (m *Migrator) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT PRIMARY KEY AUTO_INCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

// Migrate ejecuta todas las migraciones pendientes
func (m *Migrator) Migrate() error {
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("error creating database table: %v", err)
	}

	executed, err := m.getExecutedMigrations()
	if err != nil {
		return err
	}

	// Ordenar migraciones por timestamp
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].GetTimestamp() < m.migrations[j].GetTimestamp()
	})

	batch, err := m.getNextBatch()
	if err != nil {
		return err
	}

	for _, migration := range m.migrations {
		migrationName := fmt.Sprintf("%s_%s", migration.GetTimestamp(), migration.GetName())

		if _, exists := executed[migrationName]; !exists {
			fmt.Printf("Migrating: %s\n", migrationName)

			if err := migration.Up(m.db); err != nil {
				return fmt.Errorf("error executing migration %s: %v", migrationName, err)
			}

			if err := m.recordMigration(migrationName, batch); err != nil {
				return err
			}

			fmt.Printf("Migrated: %s\n", migrationName)
		}
	}

	return nil
}

func (m *Migrator) Fresh() error {
	// Eliminar la tabla de migraciones
	if err := dropAllMigrationsTable(m.db); err != nil {
		return fmt.Errorf("error dropping migrations table: %v", err)
	}

	// Volver a crear la tabla de migraciones
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("error recreating migrations table: %v", err)
	}

	// Ejecutar todas las migraciones nuevamente
	if err := m.Migrate(); err != nil {
		return fmt.Errorf("error running migrations after fresh: %v", err)
	}

	fmt.Println("Migrations table has been refreshed successfully.")
	return nil
}

// Rollback revierte el último lote de migraciones
func (m *Migrator) Rollback() error {
	lastBatch, err := m.getLastBatch()
	if err != nil {
		return err
	}

	if lastBatch == 0 {
		fmt.Println("Nothing to rollback")
		return nil
	}

	migrations, err := m.getMigrationsByBatch(lastBatch)
	if err != nil {
		return err
	}

	// Ejecutar rollback en orden inverso
	for i := len(migrations) - 1; i >= 0; i-- {
		migrationName := migrations[i]
		migration := m.findMigrationByName(migrationName)

		if migration == nil {
			return fmt.Errorf("migration %s not found in registered database", migrationName)
		}

		fmt.Printf("Rolling back: %s\n", migrationName)

		if err := migration.Down(m.db); err != nil {
			return fmt.Errorf("error rolling back migration %s: %v", migrationName, err)
		}

		if err := m.deleteMigrationRecord(migrationName); err != nil {
			return err
		}

		fmt.Printf("Rolled back: %s\n", migrationName)
	}

	return nil
}

func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT migration FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		executed[migration] = true
	}

	return executed, nil
}

func (m *Migrator) getNextBatch() (int, error) {
	var batch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM migrations").Scan(&batch)
	return batch, err
}

func (m *Migrator) getLastBatch() (int, error) {
	var batch int
	err := m.db.QueryRow("SELECT COALESCE(MAX(batch), 0) FROM migrations").Scan(&batch)
	return batch, err
}

func (m *Migrator) getMigrationsByBatch(batch int) ([]string, error) {
	rows, err := m.db.Query("SELECT migration FROM migrations WHERE batch = ? ORDER BY id DESC", batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (m *Migrator) recordMigration(name string, batch int) error {
	_, err := m.db.Exec("INSERT INTO migrations (migration, batch) VALUES (?, ?)", name, batch)
	return err
}

func (m *Migrator) deleteMigrationRecord(name string) error {
	_, err := m.db.Exec("DELETE FROM migrations WHERE migration = ?", name)
	return err
}

func (m *Migrator) findMigrationByName(name string) Migration {
	for _, migration := range m.migrations {
		migrationName := fmt.Sprintf("%s_%s", migration.GetTimestamp(), migration.GetName())
		if migrationName == name {
			return migration
		}
	}
	return nil
}

func dropAllMigrationsTable(db *sql.DB) error {
	// Deshabilitar claves foráneas
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	var rows, errorRows = db.Query("SHOW TABLES")
	if errorRows != nil {
		return errorRows
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if errorScan := rows.Scan(&tableName); errorScan != nil {
			return errorScan
		}
		tables = append(tables, tableName)
	}

	// Eliminar primero las tablas que contienen 'tokens' (hijas)
	for _, tableName := range tables {
		if tableName == "migrations" {
			continue
		}
		if containsToken(tableName) {
			_, errorExecute := db.Exec("DROP TABLE IF EXISTS " + tableName)
			if errorExecute != nil {
				return errorExecute
			}
		}
	}
	// Luego eliminar el resto
	for _, tableName := range tables {
		if tableName == "migrations" {
			continue
		}
		if !containsToken(tableName) {
			_, errorExecute := db.Exec("DROP TABLE IF EXISTS " + tableName)
			if errorExecute != nil {
				return errorExecute
			}
		}
	}
	// Finalmente, eliminar la tabla de migraciones
	_, _ = db.Exec("DROP TABLE IF EXISTS migrations")

	// Volver a habilitar claves foráneas
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	return nil
}

func containsToken(tableName string) bool {
	// Busca si la palabra 'token' está en cualquier parte del nombre de la tabla
	return len(tableName) >= 5 && (tableName == "tokens" || tableName == "token" || containsSubstring(tableName, "token"))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
