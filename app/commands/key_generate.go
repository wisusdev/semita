package commands

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"semita/app/utils"

	"github.com/spf13/cobra"
)

var KeyGenerateCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Genera una nueva clave JWT y la guarda en el archivo .env",
	Run: func(cmd *cobra.Command, args []string) {
		key := generateRandomKey(32)
		fmt.Println("Clave generada:", key)

		// Intenta actualizar el archivo .env
		utils.UpdateEnvFile("APP_KEY", key)
	},
}

func generateRandomKey(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
