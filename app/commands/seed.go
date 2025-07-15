package commands

import (
	"fmt"
	"log"
	"semita/app/core/database"
	"semita/app/utils"
	"semita/database/seeders"

	"github.com/spf13/cobra"
)

// SeedAllCommand ejecuta todos los seeders
var SeedAllCommand = &cobra.Command{
	Use:   "db:seed",
	Short: "Ejecuta todos los seeders",
	Long:  "Execute all registered seeders in the correct dependency order.",
	Run: func(cmd *cobra.Command, args []string) {
		runAllSeeders()
	},
}

// SeedRunCommand ejecuta un seeder específico
var SeedRunCommand = &cobra.Command{
	Use:   "run:seed [seeder_name]",
	Short: "Ejecuta un seeder específico",
	Long:  "Execute a specific seeder by name. Dependencies will be run first if needed.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSpecificSeeder(args[0])
	},
}

// createSeederManager crea y configura el manager de seeders
func createSeederManager() *database.SeederManager {
	manager := database.NewSeederManager()

	// Registrar todos los seeders
	manager.RegisterSeeder(seeders.NewRolesPermissionsSeeder())
	manager.RegisterSeeder(seeders.NewCategoriesSeeder())
	manager.RegisterSeeder(seeders.NewUsersSeeder())

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
		utils.Logs("ERROR", fmt.Sprintf("%v", err))
		log.Fatalf("Error running seeder '%s': %v", seederName, err)
	}

	log.Printf("=== Seeder '%s' Completed Successfully ===", seederName)
}
