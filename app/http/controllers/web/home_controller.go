package web

import (
	"fmt"
	"text/template"
	"web_utilidades/app/helpers"
	"web_utilidades/config"

	"github.com/gin-gonic/gin"
)

type Habilidades struct {
	Nombre string
}

type Data struct {
	Nombre      string
	Edad        int
	Perfil      string
	Habilidades []Habilidades
}

func HomeIndex(c *gin.Context) {
	authHomeData := helpers.AuthSessionService(c.Writer, c.Request, "Inicio", nil)
	templateHome := template.Must(template.ParseFiles("resources/home/home.html", config.MainLayoutFilePath))
	if err := templateHome.Execute(c.Writer, authHomeData); err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func Nosotros(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "Nosotros", nil)
	templateWe := template.Must(template.ParseFiles("resources/home/nosotros.html", config.MainLayoutFilePath))

	err := templateWe.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func Parametros(context *gin.Context) {
	data := map[string]string{
		"id":   context.Param("id"),
		"slug": context.Param("slug"),
	}
	data["text"] = "Hello World"
	authData := helpers.AuthSessionService(context.Writer, context.Request, "Parametros", data)
	tmpl, err := template.ParseFiles("resources/home/parametros.html", config.MainLayoutFilePath)
	if err != nil {
		fmt.Println("Error al cargar la plantilla")
		fmt.Fprintln(context.Writer, "Error al cargar la plantilla")
		return
	}
	tmpl.Execute(context.Writer, authData)
}

func QueryString(context *gin.Context) {
	data := map[string]string{
		"id":   context.Query("id"),
		"slug": context.Query("slug"),
	}
	authData := helpers.AuthSessionService(context.Writer, context.Request, "Query String", data)
	tmpl, err := template.ParseFiles("resources/home/querystring.html", config.MainLayoutFilePath)
	if err != nil {
		fmt.Println("Error al cargar la plantilla")
		fmt.Fprintln(context.Writer, "Error al cargar la plantilla")
		return
	}
	tmpl.Execute(context.Writer, authData)
}

func Estructuras(c *gin.Context) {
	habilidad01 := Habilidades{"Desarrollador Backend"}
	habilidad02 := Habilidades{"Desarrollador Frontend"}
	habilidad03 := Habilidades{"Desarrollador Fullstack"}

	data := Data{
		Nombre:      "Juan",
		Edad:        18,
		Perfil:      "Administrador de Sistemas",
		Habilidades: []Habilidades{habilidad01, habilidad02, habilidad03},
	}
	authData := helpers.AuthSessionService(c.Writer, c.Request, "Estructuras", data)
	templateFile, _ := template.ParseFiles("resources/home/estructuras.html", config.MainLayoutFilePath)
	err := templateFile.Execute(c.Writer, authData)
	if err != nil {
		return
	}
}
