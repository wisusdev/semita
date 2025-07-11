package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// GoDotEnv carga el archivo .env una sola vez al inicio de la aplicación.
func GoDotEnv() {
	err := godotenv.Load()
	if err != nil {
		Logs("error", "Error loading .env file")
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}
}

// GetEnv obtiene una variable de entorno, o retorna un valor por defecto si no existe.
func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		Logs("info", fmt.Sprintf("La variable de entorno %s no está definida.", key))
	}
	return val
}

func UpdateEnvFile(key, value string) {
	file := ".env"
	input, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("No se pudo leer %s, solo se mostrará la clave.\n", file)
		return
	}
	lines := strings.Split(string(input), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, key+"="+value)
	}
	output := strings.Join(lines, "\n")
	os.WriteFile(file, []byte(output), 0644)
	fmt.Printf("%s actualizado en %s\n", key, file)
}
