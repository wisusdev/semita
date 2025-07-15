package bootstrap

import (
	"os"
	"semita/app/commands"

	"github.com/spf13/cobra"
)

func Commands() {
	// Registra los comandos
	RootCmd.AddCommand(commands.MigrateCmd)
	RootCmd.AddCommand(commands.MigrateFreshCmd)
	RootCmd.AddCommand(commands.MigrateRollbackCmd)
	RootCmd.AddCommand(commands.MakeMigrationCmd)
	RootCmd.AddCommand(commands.MakeMigrationFromDbCmd)
	RootCmd.AddCommand(commands.KeyGenerateCmd)
	RootCmd.AddCommand(commands.OauthKeysCmd)
	RootCmd.AddCommand(commands.OauthClientCmd)
	RootCmd.AddCommand(commands.SeedAllCommand)
	RootCmd.AddCommand(commands.SeedRunCommand)

	// Ejecuta la CLI
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "semita",
	Short: "Semita CLI",
	// Si no hay subcomando, muestra la ayuda
	Run: func(cmd *cobra.Command, args []string) {
		// Si no hay subcomando, arrancar el servidor
		if len(args) == 0 && (os.Getenv("AIR") != "" || len(os.Args) == 1) {
			commands.StartServer()
			return
		}
		cmd.Help()
	},
}
