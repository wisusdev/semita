package main

import (
	"web_utilidades/app/utils"
	"web_utilidades/bootstrap"
)

func main() {
	utils.GoDotEnv()         // Cargar .env solo una vez al inicio
	utils.LoadTranslations() // Cargar traducciones
	bootstrap.Commands()
}
