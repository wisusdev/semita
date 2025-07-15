package web

import (
	"fmt"
	"net/http"
	"semita/app/helpers"
	"semita/app/models"
	"semita/app/notifications"
	"semita/app/structs"
	"semita/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthLogin(context *gin.Context) {
	helpers.View(context, "auth/login.html", "Login", nil)
}

func AuthLoginPost(context *gin.Context) {
	email := context.PostForm("email")
	password := context.PostForm("password")

	if email == "" || password == "" {
		utils.Logs("ERROR", "Email and password are required")
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Email and password are required")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	user := structs.LoginUserStruct{
		Email:    email,
		Password: password,
	}

	storedUser, err := models.GetUserByEmail(user.Email)
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Error retrieving user: %v", err))
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Invalid email or password")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if errPassword != nil {
		utils.Logs("ERROR", "Invalid password")
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Invalid email or password")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	sessionLoginError := utils.LoginUserSession(context.Writer, context.Request, storedUser)
	if sessionLoginError != nil {
		utils.Logs("ERROR", fmt.Sprintf("Error creating user session: %v", sessionLoginError))
		utils.CreateFlashNotification(context.Writer, context.Request, "error", "Error creating user session")
		context.Redirect(http.StatusSeeOther, "/auth/login")
		context.Abort()
		return
	}

	utils.CreateFlashNotification(context.Writer, context.Request, "success", "Login successful!")
	context.Redirect(http.StatusSeeOther, "/")
	context.Abort()
}

func AuthLogout(c *gin.Context) {
	sessionLogoutError := utils.LogoutUserSession(c.Writer, c.Request)
	if sessionLogoutError != nil {
		c.String(http.StatusInternalServerError, "Error logging out")
		return
	}

	utils.CreateFlashNotification(c.Writer, c.Request, "success", "Logout successful!")
	c.Redirect(http.StatusSeeOther, "/")
	c.Abort()
}

func AuthRegister(context *gin.Context) {
	helpers.View(context, "auth/register.html", "Register", nil)
}

func AuthRegisterPost(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	if name == "" || email == "" || password == "" || confirmPassword == "" {
		c.String(http.StatusBadRequest, "Name, Email, password, and confirm password are required")
		return
	}

	if password != confirmPassword {
		c.String(http.StatusBadRequest, "Passwords do not match")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error encrypting password")
		return
	}

	user := structs.StoreUserStruct{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	errorStore := models.StoreUser(user)
	if errorStore != nil {
		c.String(http.StatusInternalServerError, "Error saving user to the database")
		return
	}

	c.Redirect(http.StatusSeeOther, "/auth/login")
	c.Abort()
}

func AuthForgotPassword(context *gin.Context) {
	helpers.View(context, "auth/forgot_password.html", "Recuperar Contraseña", nil)
}

func AuthForgotPasswordPost(context *gin.Context) {
	email := context.PostForm("email")
	if email == "" {
		context.String(http.StatusBadRequest, "Email is required")
		return
	}

	token := utils.GenerateResetToken(email)
	resetURL := "http://" + utils.GetEnv("APP_URL") + "/auth/reset-password?token=" + token
	_ = models.CreatePasswordReset(email, token) // Guardar token en BD
	errorSendEmail := notifications.SendPasswordReset(email, resetURL)

	if errorSendEmail != nil {
		utils.Logs("ERROR", errorSendEmail.Error())
		fmt.Println("Error sending password reset email:", errorSendEmail)
		return
	}

	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}

func AuthResetPassword(context *gin.Context) {
	var data = map[string]string{
		"token": context.Query("token"),
	}

	helpers.View(context, "auth/reset_password.html", "Restablecer Contraseña", data)
}

func AuthResetPasswordPost(context *gin.Context) {
	token := context.PostForm("token")
	password := context.PostForm("password")
	confirmPassword := context.PostForm("confirm_password")

	if token == "" || password == "" || confirmPassword == "" {
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Token, password, and confirm password are required")
		return
	}

	if password != confirmPassword {
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Passwords do not match")
		return
	}

	passwordResetByToken, err := models.GetPasswordResetByToken(token)
	if err != nil {
		utils.Logs("ERROR", err.Error())
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Token inválido o expirado")
		context.Redirect(http.StatusSeeOther, "/auth/reset-password?token="+token)
		context.Abort()
		return
	}

	// Usar hora local para ambos tiempos
	now := time.Now()
	tokenCreatedAt := passwordResetByToken.CreatedAt
	timeSince := now.Sub(tokenCreatedAt)

	// Verificar expiración de 2 horas
	if timeSince > 2*time.Hour {
		_ = models.DeletePasswordReset(token)
		utils.Logs("INFO", fmt.Sprintf("Token expirado. Creado hace: %v", timeSince))
		utils.CreateFlashNotification(context.Writer, context.Request, "error", "Token expirado. Por favor, solicita un nuevo enlace de restablecimiento.")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	user, err := models.GetUserByEmail(passwordResetByToken.Email)
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Usuario no encontrado: %v", err))
		utils.CreateFlashNotification(context.Writer, context.Request, "warning", "Usuario no encontrado")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("Error al encriptar contraseña: %v", err))
		utils.CreateFlashNotification(context.Writer, context.Request, "error", "Error al encriptar contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	update := structs.UpdateUserStruct{ID: user.ID, Name: user.Name, Email: user.Email, Password: string(hashedPassword)}
	err = models.UpdateUser(update)

	if err != nil {
		utils.Logs("ERROR", fmt.Sprintf("No se pudo actualizar la contraseña: %v", err))
		utils.CreateFlashNotification(context.Writer, context.Request, "error", "No se pudo actualizar la contraseña")
		context.Redirect(http.StatusSeeOther, "/auth/forgot-password")
		context.Abort()
		return
	}

	// Eliminar el token después de usarlo exitosamente
	_ = models.DeletePasswordReset(token)
	utils.Logs("INFO", "Contraseña restablecida exitosamente")

	utils.CreateFlashNotification(context.Writer, context.Request, "success", "Contraseña actualizada exitosamente!")
	context.Redirect(http.StatusSeeOther, "/auth/login")
	context.Abort()
}
