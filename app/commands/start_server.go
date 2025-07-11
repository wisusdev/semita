package commands

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"web_utilidades/app/http/controllers/web"
	"web_utilidades/app/utils"
	"web_utilidades/routes"
)

func StartServer() {
	// Cargar variables de entorno
	var appUrl = utils.GetEnv("APP_URL")

	// Inicializar el enrutador Gin
	router := routes.Web()

	// Montar rutas API
	apiGroup := router.Group("/api/v1")
	routes.Api(apiGroup)

	// Archivos estáticos
	router.Static("/public", "./public")

	// Ruta 404 personalizada
	router.NoRoute(web.Error404)

	// Ejecución del servidor
	server := &http.Server{
		Addr:         appUrl,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Servidor corriendo en http://%v a las %s\n", appUrl, time.Now().Format("2006-01-02 15:04:05"))
	log.Fatal(server.ListenAndServe())
}
