package middleware

import (
	"net/http"
	"semita/app/utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !utils.IsUserAuthenticated(context.Request) {
			utils.CreateFlashNotification(context.Writer, context.Request, "error", "You must be logged in to access this page.")
			context.Redirect(http.StatusSeeOther, "/auth/login")
			context.Abort()
			return
		}
		handler(context)
	}
}

func RedirectGuest(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if utils.IsUserAuthenticated(context.Request) {
			context.Redirect(http.StatusSeeOther, "/")
			context.Abort()
			return
		}
		handler(context)
	}
}
