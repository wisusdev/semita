package commands

import (
	"log"
	"web_utilidades/app/core/database"
	"web_utilidades/database/seeders"

	"github.com/spf13/cobra"
)

// SeedCommand comando principal para seeders
var SeedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Database seeding commands",
	Long:  "Commands to manage database seeders including running, rolling back, and checking status.",
}

// SeedAllCommand ejecuta todos los seeders
var SeedAllCommand = &cobra.Command{
	Use:   "all",
	Short: "Run all seeders",
	Long:  "Execute all registered seeders in the correct dependency order.",
	Run: func(cmd *cobra.Command, args []string) {
		runAllSeeders()
	},
}

// SeedRunCommand ejecuta un seeder específico
var SeedRunCommand = &cobra.Command{
	Use:   "run [seeder_name]",
	Short: "Run a specific seeder",
	Long:  "Execute a specific seeder by name. Dependencies will be run first if needed.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSpecificSeeder(args[0])
	},
}

// SeedRollbackCommand revierte un seeder específico
var SeedRollbackCommand = &cobra.Command{
	Use:   "rollback [seeder_name]",
	Short: "Rollback a specific seeder",
	Long:  "Rollback (undo) a specific seeder by name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rollbackSpecificSeeder(args[0])
	},
}

// SeedResetCommand reinicia un seeder (rollback + seed)
var SeedResetCommand = &cobra.Command{
	Use:   "reset [seeder_name]",
	Short: "Reset a specific seeder",
	Long:  "Reset a specific seeder (rollback then seed again).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resetSpecificSeeder(args[0])
	},
}

// SeedStatusCommand muestra el estado de todos los seeders
var SeedStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show seeder status",
	Long:  "Display the execution status of all registered seeders.",
	Run: func(cmd *cobra.Command, args []string) {
		showSeederStatus()
	},
}

// SeedCleanCommand limpia todos los datos de seeders
var SeedCleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Clean all seeder data",
	Long:  "Remove all data created by seeders (runs rollback on all seeders).",
	Run: func(cmd *cobra.Command, args []string) {
		cleanAllSeederData()
	},
}

// createSeederManager crea y configura el manager de seeders
func createSeederManager() *database.SeederManager {
	manager := database.NewSeederManager()

	// Registrar todos los seeders
	manager.RegisterSeeder(seeders.NewRolesPermissionsSeeder())
	manager.RegisterSeeder(seeders.NewCategoriesSeeder())
	manager.RegisterSeeder(seeders.NewUsersSeeder())
	manager.RegisterSeeder(seeders.NewPostsSeeder())

	return manager
}

// runAllSeeders ejecuta todos los seeders
func runAllSeeders() {
	log.Println("=== Running All Seeders ===")

	manager := createSeederManager()
	err := manager.RunAllSeeders()
	if err != nil {
		log.Fatalf("Error running all seeders: %v", err)
	}

	log.Println("=== All Seeders Completed Successfully ===")
}

// runSpecificSeeder ejecuta un seeder específico
func runSpecificSeeder(seederName string) {
	log.Printf("=== Running Seeder: %s ===", seederName)

	manager := createSeederManager()
	err := manager.RunSeeder(seederName)
	if err != nil {
		log.Fatalf("Error running seeder '%s': %v", seederName, err)
	}

	log.Printf("=== Seeder '%s' Completed Successfully ===", seederName)
}

// rollbackSpecificSeeder revierte un seeder específico
func rollbackSpecificSeeder(seederName string) {
	log.Printf("=== Rolling Back Seeder: %s ===", seederName)

	manager := createSeederManager()
	err := manager.RollbackSeeder(seederName)
	if err != nil {
		log.Fatalf("Error rolling back seeder '%s': %v", seederName, err)
	}

	log.Printf("=== Seeder '%s' Rolled Back Successfully ===", seederName)
}

// resetSpecificSeeder reinicia un seeder específico
func resetSpecificSeeder(seederName string) {
	log.Printf("=== Resetting Seeder: %s ===", seederName)

	manager := createSeederManager()
	err := manager.ResetSeeder(seederName)
	if err != nil {
		log.Fatalf("Error resetting seeder '%s': %v", seederName, err)
	}

	log.Printf("=== Seeder '%s' Reset Successfully ===", seederName)
}

// showSeederStatus muestra el estado de todos los seeders
func showSeederStatus() {
	log.Println("=== Seeder Status ===")

	manager := createSeederManager()
	manager.GetSeederStatus()
}

// cleanAllSeederData limpia todos los datos de seeders
func cleanAllSeederData() {
	log.Println("=== Cleaning All Seeder Data ===")

	manager := createSeederManager()
	err := manager.CleanAllSeederData()
	if err != nil {
		log.Fatalf("Error cleaning seeder data: %v", err)
	}

	log.Println("=== All Seeder Data Cleaned Successfully ===")
}

func init() {
	// Agregar subcomandos al comando principal
	SeedCommand.AddCommand(SeedAllCommand)
	SeedCommand.AddCommand(SeedRunCommand)
	SeedCommand.AddCommand(SeedRollbackCommand)
	SeedCommand.AddCommand(SeedResetCommand)
	SeedCommand.AddCommand(SeedStatusCommand)
	SeedCommand.AddCommand(SeedCleanCommand)
}
