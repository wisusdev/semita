package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"web_utilidades/app/helpers"
	"web_utilidades/app/structs"
	"web_utilidades/config"

	"github.com/gin-gonic/gin"
)

var baseUri = "http://localhost:3000/"

func DummyApiIndex(context *gin.Context) {

	var users, errorUsers = getAllUsers()

	if errorUsers != nil {
		fmt.Println("Error al obtener los usuarios:", errorUsers)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		return
	}

	helpers.View(context, "dummyjson/index.html", "API Index", users)
}

func DummyApiCreate(c *gin.Context) {
	var viewData = helpers.AuthSessionService(c.Writer, c.Request, "API create", nil)
	var templateCreate = template.Must(template.ParseFiles("resources/dummyjson/create.html", config.MainLayoutFilePath))
	var errorExecuteTemplate = templateCreate.Execute(c.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al mostrar el formulario de creación"})
		return
	}
}

func DummyApiStore(context *gin.Context) {
	// Añade log para depuración
	fmt.Println("DummyApiStore - Method:", context.Request.Method)
	if context.Request.Method == "POST" {
		var name = context.PostForm("name")
		var email = context.PostForm("email")
		var username = context.PostForm("username")

		if name == "" || email == "" || username == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		var uri = baseUri + "users"

		var header = map[string]string{
			"Content-Type": "application/json",
		}

		var body = map[string]string{
			"name":     name,
			"email":    email,
			"username": username,
		}

		var responseData, errorResponse = helpers.MakeRequest("POST", uri, body, header, true)
		if errorResponse != nil {
			fmt.Println("Error al crear el usuario:", errorResponse)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
			return
		}

		var createdUser map[string]any
		var errorUnmarshal = json.Unmarshal([]byte(responseData["body"].(string)), &createdUser)
		if errorUnmarshal != nil {
			fmt.Println("Error al deserializar la respuesta:", errorUnmarshal)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del usuario creado"})
			return
		}

		fmt.Println("Usuario creado:", createdUser)
		context.Redirect(http.StatusSeeOther, "/dummyjson")
	}
}

func DummyApiShow(context *gin.Context) {
	var userId = context.Param("id")

	if userId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var uri = baseUri + "users/" + userId
	var header = map[string]string{
		"Content-Type": "application/json",
	}
	var responseData, errorResponse = helpers.MakeRequest("GET", uri, nil, header, false)
	if errorResponse != nil {
		fmt.Println("Error al obtener el usuario:", errorResponse)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el usuario"})
		return
	}
	var user structs.JsonApiUsers
	var errorUnmarshal = json.Unmarshal([]byte(responseData["body"].(string)), &user)
	if errorUnmarshal != nil {
		fmt.Println("Error al deserializar la respuesta:", errorUnmarshal)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del usuario"})
		return
	}

	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "API show", user)

	var templateShow = template.Must(template.ParseFiles("resources/dummyjson/show.html", config.MainLayoutFilePath))
	var errorExecuteTemplate = templateShow.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al mostrar el usuario"})
		return
	}
}

func DummyApiEdit(context *gin.Context) {
	var userId = context.Param("id")

	if userId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var uri = baseUri + "users/" + userId
	var header = map[string]string{
		"Content-Type": "application/json",
	}
	var responseData, errorResponse = helpers.MakeRequest("GET", uri, nil, header, false)
	if errorResponse != nil {
		fmt.Println("Error al obtener el usuario:", errorResponse)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el usuario"})
		return
	}

	var user structs.JsonApiUsers
	var errorUnmarshal = json.Unmarshal([]byte(responseData["body"].(string)), &user)
	if errorUnmarshal != nil {
		fmt.Println("Error al deserializar la respuesta:", errorUnmarshal)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del usuario"})
		return
	}

	var viewData = helpers.AuthSessionService(context.Writer, context.Request, "API Edit", user)

	var templateEdit = template.Must(template.ParseFiles("resources/dummyjson/edit.html", config.MainLayoutFilePath))
	var errorExecuteTemplate = templateEdit.Execute(context.Writer, viewData)
	if errorExecuteTemplate != nil {
		fmt.Println("Error al ejecutar la plantilla:", errorExecuteTemplate)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al mostrar el formulario de edición"})
		return
	}
}

func DummyApiUpdate(context *gin.Context) {
	if context.Request.Method != "PUT" {
		fmt.Printf("Método esperado PUT, recibido: %s\n", context.Request.Method)
		context.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método no permitido"})
		return
	}

	var userId = context.Param("id")
	if userId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Para formularios con method override, los datos están en context.Request.Form
	var name = context.Request.FormValue("name")
	var email = context.Request.FormValue("email")
	var username = context.Request.FormValue("username")

	// Si FormValue no funciona, intenta con PostForm
	if name == "" {
		name = context.PostForm("name")
	}
	if email == "" {
		email = context.PostForm("email")
	}
	if username == "" {
		username = context.PostForm("username")
	}

	fmt.Printf("Datos recibidos - Name: %s, Email: %s, Username: %s\n", name, email, username)

	if name == "" || email == "" || username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	var uri = baseUri + "users/" + userId

	var header = map[string]string{
		"Content-Type": "application/json",
	}

	var body = map[string]string{
		"name":     name,
		"email":    email,
		"username": username,
	}

	var responseData, errorResponse = helpers.MakeRequest("PATCH", uri, body, header, true)
	if errorResponse != nil {
		fmt.Println("Error al actualizar el usuario:", errorResponse)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	var updatedUser map[string]any
	var errorUnmarshal = json.Unmarshal([]byte(responseData["body"].(string)), &updatedUser)
	if errorUnmarshal != nil {
		fmt.Println("Error al deserializar la respuesta:", errorUnmarshal)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la respuesta del usuario actualizado"})
		return
	}

	context.Redirect(http.StatusSeeOther, "/dummyjson")
}

func DummyApiDelete(context *gin.Context) {
	// Verifica que el método sea DELETE después del middleware
	if context.Request.Method != "DELETE" {
		fmt.Printf("Método esperado DELETE, recibido: %s\n", context.Request.Method)
		context.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método no permitido"})
		return
	}

	var userId = context.Param("id")
	if userId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var uri = baseUri + "users/" + userId
	var header = map[string]string{
		"Content-Type": "application/json",
	}

	var _, errorResponse = helpers.MakeRequest("DELETE", uri, nil, header, false)
	if errorResponse != nil {
		fmt.Println("Error al eliminar el usuario:", errorResponse)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el usuario"})
		return
	}

	context.Redirect(http.StatusSeeOther, "/dummyjson")
}

func getAllUsers() (structs.JsonApiStruct, error) {

	var uri = baseUri + "users"
	var header = map[string]string{
		"Content-Type": "application/json",
	}

	var response, errorResponse = helpers.MakeRequest("GET", uri, nil, header, false)
	if errorResponse != nil {
		fmt.Println("Error al obtener los usuarios:", errorResponse)
		return structs.JsonApiStruct{}, errorResponse
	}

	var userResponse structs.JsonApiStruct
	var errorUnmarshal error = json.Unmarshal([]byte(response["body"].(string)), &userResponse)
	if errorUnmarshal != nil {
		fmt.Println("Error al deserializar la respuesta:", errorUnmarshal)
		return structs.JsonApiStruct{}, errorUnmarshal
	}

	return userResponse, nil
}
