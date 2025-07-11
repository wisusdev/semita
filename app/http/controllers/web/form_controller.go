package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
	"web_utilidades/app/helpers"
	"web_utilidades/app/utils"
	validaciones "web_utilidades/app/validations"
	"web_utilidades/config"

	"github.com/gin-gonic/gin"
)

func FormulariosGet(c *gin.Context) {
	authData := helpers.AuthSessionService(c.Writer, c.Request, "Formulario", nil)
	tmpl := template.Must(template.ParseFiles("resources/formulario/form.html", config.MainLayoutFilePath))
	err := tmpl.Execute(c.Writer, authData)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		return
	}
}

func FormulariosPost(context *gin.Context) {
	/*var nombre string = request.FormValue("nombre")
	var email string = request.FormValue("email")
	var telefono string = request.FormValue("telefono")
	var password string = request.FormValue("password")

	var data map[string]string = map[string]string{
		"nombre":   nombre,
		"email":    email,
		"telefono": telefono,
		"password": password,
	}

	fmt.Println(data)*/

	// Parsear el formulario
	errorParseForm := context.Request.ParseMultipartForm(10 << 20) // 10 MB
	if errorParseForm != nil {
		http.Error(context.Writer, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Validar los datos
	mensaje := ""

	if len(context.Request.FormValue("nombre")) == 0 {
		mensaje = mensaje + " El campo nombre está vacío."
	}

	if len(context.Request.FormValue("email")) == 0 {
		mensaje = mensaje + " El campo E-Mail está vacío."
	}

	if validaciones.RegexCorreo.FindStringSubmatch(context.Request.FormValue("email")) == nil {
		mensaje = mensaje + " . El E-Mail ingresado no es válido "
	}

	if validaciones.ValidarPassword(context.Request.FormValue("password")) == false {
		mensaje = mensaje + " . La contraseña debe tener al menos 1 número, una mayúscula, y un largo entre 6 y 8 caracteres "
	}

	// Obtenemos el archivo
	var file, handler, errorFormFile = context.Request.FormFile("archivo")
	if errorFormFile != nil {
		mensaje = mensaje + " . No se ha subido ningún archivo. "

		utils.CreateFlashNotification(context.Writer, context.Request, "warning", mensaje)
		context.Redirect(http.StatusSeeOther, "/formulario")
		return
	}

	// Validamos si el achivo es diferente de nil
	if file != nil {
		defer file.Close()
		var extension = strings.Split(handler.Filename, ".")[1]
		var timeNow = time.Now().Format("20060102150405")
		var fileUpload = timeNow + "." + extension
		var filePath = "storage/files/" + fileUpload

		var fileDestination, errorCreateFile = os.Create(filePath)
		if errorCreateFile != nil {
			utils.CreateFlashNotification(context.Writer, context.Request, "danger", "Error al crear el archivo")
			context.Redirect(http.StatusSeeOther, "/formulario")
			return
		}
		defer fileDestination.Close()

		_, errorCopyFile := io.Copy(fileDestination, file)
		if errorCopyFile != nil {
			utils.CreateFlashNotification(context.Writer, context.Request, "danger", "Error al copiar el archivo")
			context.Redirect(http.StatusSeeOther, "/formulario")
			return
		}

		utils.CreateFlashNotification(context.Writer, context.Request, "success", "Archivo subido correctamente: "+fileUpload)
		context.Redirect(http.StatusSeeOther, "/formulario")
		return
	}

	if mensaje != "" {
		fmt.Println(context.Writer, mensaje)
		utils.CreateFlashNotification(context.Writer, context.Request, "danger", mensaje)
		context.Redirect(http.StatusSeeOther, "/formulario")
	}
}
