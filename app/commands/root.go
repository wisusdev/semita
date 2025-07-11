package commands

import (
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia el servidor web",
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}
