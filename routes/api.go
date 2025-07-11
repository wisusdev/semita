package routes

import (
	"github.com/gin-gonic/gin"
	"web_utilidades/app/http/controllers/api/v1/auth"
	"web_utilidades/app/http/middleware"
)

func Api(router *gin.RouterGroup) {
	// Auth routes
	router.POST("/auth/login", auth.Login)
	router.POST("/auth/register", auth.Register)
	router.POST("/auth/logout", middleware.AuthMiddleware(), auth.Logout)
	router.POST("/auth/forgot-password", auth.ForgotPassword)
	router.POST("/auth/reset-password", auth.ResetPassword)
	router.POST("/auth/email/resend", middleware.AuthMiddleware(), auth.ResendEmailVerify)
	router.GET("/auth/email/verify/:id/:hash", auth.VerifyEmail)
	router.POST("/auth/refresh-token", middleware.AuthMiddleware(), auth.RefreshToken)
}
