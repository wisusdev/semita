package routes

import (
	"web_utilidades/app/http/controllers/web"
	"web_utilidades/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func Web() *gin.Engine {
	router := gin.Default()

	// IMPORTANTE: El middleware debe estar ANTES de todas las rutas
	router.Use(middleware.MethodOverride(), middleware.LanguageMiddleware())

	// Ahora define todas las rutas
	router.GET("/", web.HomeIndex)

	// Auth routes
	router.GET("/auth/login", middleware.RedirectGuest(web.AuthLogin))
	router.POST("/auth/login", middleware.RedirectGuest(web.AuthLoginPost))
	router.GET("/auth/logout", middleware.RequireAuth(web.AuthLogout))
	router.GET("/auth/register", middleware.RedirectGuest(web.AuthRegister))
	router.POST("/auth/register", middleware.RedirectGuest(web.AuthRegisterPost))
	router.GET("/auth/forgot-password", web.AuthForgotPassword)
	router.POST("/auth/forgot-password", web.AuthForgotPasswordPost)
	router.GET("/auth/reset-password", web.AuthResetPassword)
	router.POST("/auth/reset-password", web.AuthResetPasswordPost)

	// General routes
	router.GET("/nosotros", middleware.RequireAuth(web.Nosotros))
	router.GET("/parametros/:id/:slug", middleware.RequireAuth(web.Parametros))
	router.GET("/querystring", middleware.RequireAuth(web.QueryString))
	router.GET("/estructuras", middleware.RequireAuth(web.Estructuras))

	// Form routes
	router.GET("/formulario", middleware.RequireAuth(web.FormulariosGet))
	router.POST("/formulario-post", middleware.RequireAuth(web.FormulariosPost))

	// Utility routes
	router.GET("/pdf", middleware.RequireAuth(web.IndexPDF))
	router.GET("/pdf/new", middleware.RequireAuth(web.GenerateNewPDF))
	router.GET("/excel", middleware.RequireAuth(web.IndexExcel))
	router.GET("/excel/new", middleware.RequireAuth(web.GenerateNewExcel))
	router.GET("/qr", middleware.RequireAuth(web.IndexQR))
	router.GET("/qr/new", middleware.RequireAuth(web.GenerateNewQR))
	router.GET("/email", middleware.RequireAuth(web.IndexSendEmail))
	router.GET("/email/new", middleware.RequireAuth(web.GenerateNewEmail))

	router.GET("/dummyjson", middleware.RequireAuth(web.DummyApiIndex))
	router.GET("/dummyjson/users/create", middleware.RequireAuth(web.DummyApiCreate))
	router.POST("/dummyjson/users/store", middleware.RequireAuth(web.DummyApiStore))
	router.GET("/dummyjson/users/show/:id", middleware.RequireAuth(web.DummyApiShow))
	router.GET("/dummyjson/users/edit/:id", middleware.RequireAuth(web.DummyApiEdit))
	router.POST("/dummyjson/users/update/:id", middleware.RequireAuth(web.DummyApiUpdate))
	router.POST("/dummyjson/users/delete/:id", middleware.RequireAuth(web.DummyApiDelete))

	router.GET("/users", middleware.RequireAuth(web.UserIndex))
	router.GET("/users/create", middleware.RequireAuth(web.UserCreate))
	router.POST("/users/store", middleware.RequireAuth(web.UserStore))
	router.GET("/users/show/:id", middleware.RequireAuth(web.UserShow))
	router.GET("/users/edit/:id", middleware.RequireAuth(web.UserEdit))
	router.POST("/users/update/:id", middleware.RequireAuth(web.UserUpdate))
	router.POST("/users/delete/:id", middleware.RequireAuth(web.UserDelete))

	// Ruta para cambiar idioma
	router.POST("/set-lang", middleware.SetLangHandler)

	return router
}
