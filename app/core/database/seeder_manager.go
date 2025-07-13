package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"web_utilidades/app/utils"
	"web_utilidades/config"
)

// Seeder interface que deben implementar todos los seeders
type Seeder interface {
	Seed() error
	Rollback() error
	GetName() string
	GetDependencies() []string
}

// BaseSeeder estructura base para todos los seeders
type BaseSeeder struct {
	DB   *sql.DB
	Name string
}

// SeederManager gestiona la ejecución de seeders
type SeederManager struct {
	DB      *sql.DB
	seeders map[string]Seeder
}

// NewSeederManager crea una nueva instancia del manager
func NewSeederManager() *SeederManager {
	db := config.DatabaseConnect()
	return &SeederManager{
		DB:      db,
		seeders: make(map[string]Seeder),
	}
}

// RegisterSeeder registra un seeder en el manager
func (sm *SeederManager) RegisterSeeder(seeder Seeder) {
	sm.seeders[seeder.GetName()] = seeder
	utils.Logs("INFO", fmt.Sprintf("Seeder '%s' registered successfully", seeder.GetName()))
}

// GetAllSeeders retorna todos los seeders registrados
func (sm *SeederManager) GetAllSeeders() map[string]Seeder {
	return sm.seeders
}

// GetSeeder retorna un seeder específico por nombre
func (sm *SeederManager) GetSeeder(name string) (Seeder, error) {
	seeder, exists := sm.seeders[name]
	if !exists {
		utils.Logs("ERROR", fmt.Sprintf("El seeder '%s' no existe", name))
		return nil, fmt.Errorf("seeder '%s' not found", name)
	}
	return seeder, nil
}

// IsSeederExecuted siempre retorna false para forzar la re-ejecución
// con limpieza de datos cada vez
func (sm *SeederManager) IsSeederExecuted(name string) bool {
	// Siempre retornar false para que los seeders se ejecuten cada vez
	// limpiando y volviendo a poblar las tablas
	return false
}

// MarkSeederAsExecuted registra la ejecución del seeder pero no afecta
// la lógica de re-ejecución (opcional para logs/auditoría)
func (sm *SeederManager) MarkSeederAsExecuted(name string) error {
	// Opcional: mantener registro para auditoría, pero no afecta la lógica
	// Primero intentar actualizar, si no existe insertar
	updateQuery := `UPDATE seeders SET executed_at = ?, rollback_at = NULL WHERE name = ?`
	result, err := sm.DB.Exec(updateQuery, time.Now(), name)
	if err != nil {
		log.Printf("Warning: could not update seeder execution log: %v", err)
		return nil
	}

	// Verificar si se actualizó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		// No se actualizó ninguna fila, insertar nueva
		insertQuery := `INSERT INTO seeders (name, executed_at) VALUES (?, ?)`
		_, err = sm.DB.Exec(insertQuery, name, time.Now())
		if err != nil {
			log.Printf("Warning: could not log seeder execution: %v", err)
			// No retornar error ya que esto es solo para auditoría
		} else {
			log.Printf("Seeder '%s' execution logged", name)
		}
	} else {
		log.Printf("Seeder '%s' execution updated", name)
	}

	return nil
}

// MarkSeederAsRolledBack registra el rollback del seeder (opcional para auditoría)
func (sm *SeederManager) MarkSeederAsRolledBack(name string) error {
	// Opcional: mantener registro para auditoría
	query := `UPDATE seeders SET rollback_at = ? WHERE name = ?`
	_, err := sm.DB.Exec(query, time.Now(), name)
	if err != nil {
		log.Printf("Warning: could not log seeder rollback: %v", err)
		// No retornar error ya que esto es solo para auditoría
	} else {
		log.Printf("Seeder '%s' rollback logged", name)
	}
	return nil
}

// RunSeeder ejecuta un seeder específico limpiando primero los datos
func (sm *SeederManager) RunSeeder(name string) error {

	seeder, err := sm.GetSeeder(name)
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Seeder '%s' not found: %v", name, err))
		return err
	}

	// Ejecutar dependencias primero
	dependencies := seeder.GetDependencies()
	for _, dep := range dependencies {
		log.Printf("Running dependency seeder: %s", dep)
		err := sm.RunSeeder(dep)
		if err != nil {
			utils.Logs("ERROR", fmt.Sprintf("Error running dependency '%s': %v", dep, err))
			return fmt.Errorf("error running dependency '%s': %v", dep, err)
		}
	}

	// Ejecutar el seeder en una transacción
	tx, err := sm.DB.Begin()
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Error starting transaction for seeder '%s': %v", name, err))
		return fmt.Errorf("error starting transaction: %v", err)
	}

	log.Printf("Running seeder (with data cleanup): %s", name)

	// Primero hacer rollback para limpiar datos existentes
	log.Printf("Cleaning existing data for seeder: %s", name)
	err = seeder.Rollback()
	if err != nil {
		log.Printf("Warning: error during cleanup for seeder '%s': %v", name, err)
		// Continuar aunque falle la limpieza (puede no haber datos previos)
	}

	// Ahora ejecutar el seeding
	err = seeder.Seed()
	if err != nil {
		tx.Rollback()
		utils.Logs("ERROR", fmt.Sprintf("Error running seeder '%s': %v", name, err))
		return fmt.Errorf("error running seeder '%s': %v", name, err)
	}

	// Marcar como ejecutado (opcional para auditoría)
	sm.MarkSeederAsExecuted(name)

	err = tx.Commit()
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Error committing transaction for seeder '%s': %v", name, err))
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Printf("Seeder '%s' executed successfully", name)
	return nil
}

// RollbackSeeder revierte un seeder específico
func (sm *SeederManager) RollbackSeeder(name string) error {
	seeder, err := sm.GetSeeder(name)
	if err != nil {
		return err
	}

	// Ejecutar rollback en una transacción
	tx, err := sm.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	log.Printf("Rolling back seeder: %s", name)
	err = seeder.Rollback()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error rolling back seeder '%s': %v", name, err)
	}

	// Marcar como revertido (opcional para auditoría)
	sm.MarkSeederAsRolledBack(name)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Printf("Seeder '%s' rolled back successfully", name)
	return nil
}

// RunAllSeeders ejecuta todos los seeders registrados con limpieza automática
func (sm *SeederManager) RunAllSeeders() error {
	utils.Logs("INFO", "Running all seeders with data cleanup...")

	// Crear un grafo de dependencias y ejecutar en orden
	executed := make(map[string]bool)

	var runSeederWithDeps func(string) error
	runSeederWithDeps = func(name string) error {
		if executed[name] {
			return nil
		}

		seeder, exists := sm.seeders[name]
		if !exists {
			return fmt.Errorf("seeder '%s' not found", name)
		}

		// Ejecutar dependencias primero
		for _, dep := range seeder.GetDependencies() {
			err := runSeederWithDeps(dep)
			if err != nil {
				return err
			}
		}

		// Ejecutar el seeder (que ahora incluye limpieza automática)
		err := sm.RunSeeder(name)
		if err != nil {
			utils.Logs("ERROR", fmt.Sprintf("Seeder '%s' failed: %v", name, err))
			return err
		}

		executed[name] = true
		return nil
	}

	// Ejecutar todos los seeders
	for name := range sm.seeders {
		err := runSeederWithDeps(name)
		if err != nil {
			return err
		}
	}

	log.Println("All seeders executed successfully with fresh data")
	return nil
}

// GetSeederStatus retorna el estado de todos los seeders
func (sm *SeederManager) GetSeederStatus() {
	log.Println("=== Seeder Status ===")

	for name := range sm.seeders {
		status := "Not executed"
		var executedAt, rollbackAt sql.NullTime

		query := `SELECT executed_at, rollback_at FROM seeders WHERE name = ?`
		err := sm.DB.QueryRow(query, name).Scan(&executedAt, &rollbackAt)

		if err == nil {
			if rollbackAt.Valid {
				status = fmt.Sprintf("Rolled back at %s", rollbackAt.Time.Format("2006-01-02 15:04:05"))
			} else if executedAt.Valid {
				status = fmt.Sprintf("Executed at %s", executedAt.Time.Format("2006-01-02 15:04:05"))
			}
		}

		log.Printf("%-30s: %s", name, status)
	}
}

// ResetSeeder ejecuta rollback y luego seed (ahora redundante pero mantenido para compatibilidad)
func (sm *SeederManager) ResetSeeder(name string) error {
	log.Printf("Resetting seeder (cleanup + seed): %s", name)

	// Con el nuevo comportamiento, RunSeeder ya hace cleanup automáticamente
	// pero mantenemos este método para compatibilidad
	err := sm.RunSeeder(name)
	if err != nil {
		return err
	}

	log.Printf("Seeder '%s' reset successfully", name)
	return nil
}

// CleanAllSeederData limpia todos los datos de los seeders registrados
func (sm *SeederManager) CleanAllSeederData() error {
	log.Println("=== Cleaning All Seeder Data ===")

	// Crear orden inverso para rollback (considerando dependencias)
	var seederOrder []string
	processed := make(map[string]bool)

	var addSeederToOrder func(string)
	addSeederToOrder = func(name string) {
		if processed[name] {
			return
		}

		seeder, exists := sm.seeders[name]
		if !exists {
			return
		}

		// Procesar dependencias primero
		for _, dep := range seeder.GetDependencies() {
			addSeederToOrder(dep)
		}

		seederOrder = append(seederOrder, name)
		processed[name] = true
	}

	// Agregar todos los seeders en orden
	for name := range sm.seeders {
		addSeederToOrder(name)
	}

	// Ejecutar rollback en orden inverso
	for i := len(seederOrder) - 1; i >= 0; i-- {
		seederName := seederOrder[i]
		seeder := sm.seeders[seederName]

		log.Printf("Cleaning data for: %s", seederName)
		err := seeder.Rollback()
		if err != nil {
			log.Printf("Warning: error cleaning data for '%s': %v", seederName, err)
			// Continuar con los demás aunque falle uno
		}
	}

	log.Println("=== All Seeder Data Cleaned ===")
	return nil
}
