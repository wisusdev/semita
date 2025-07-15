package main

import (
	"semita/app/utils"
	"semita/bootstrap"
)

func main() {
	utils.GoDotEnv()         // Cargar .env solo una vez al inicio
	utils.LoadTranslations() // Cargar traducciones
	bootstrap.Commands()
}
