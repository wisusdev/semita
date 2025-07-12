package middleware

import (
	"net/http"
	"web_utilidades/app/models"
	"web_utilidades/app/utils"

	"github.com/gin-gonic/gin"
)

// getUserFromSession obtiene el usuario autenticado de la sesión
func getUserFromSession(c *gin.Context) (int, bool) {
	user, authenticated := utils.GetAuthenticatedUser(c.Request)
	if !authenticated {
		return 0, false
	}
	return user.ID, true
}

// RequireRole middleware que verifica si el usuario tiene un rol específico
func RequireRole(roleName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el rol
		hasRole, err := models.UserHasRoleByName(userID, roleName, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasRole {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que verifica si el usuario tiene al menos uno de los roles especificados
func RequireAnyRole(roleNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene al menos uno de los roles
		hasAnyRole, err := models.UserHasAnyRole(userID, roleNames, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAnyRole {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllRoles middleware que verifica si el usuario tiene todos los roles especificados
func RequireAllRoles(roleNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene todos los roles
		hasAllRoles, err := models.UserHasAllRoles(userID, roleNames, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAllRoles {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission middleware que verifica si el usuario tiene un permiso específico
func RequirePermission(permissionName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el permiso
		hasPermission, err := models.UserHasPermission(userID, permissionName, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission middleware que verifica si el usuario tiene al menos uno de los permisos especificados
func RequireAnyPermission(permissionNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene al menos uno de los permisos
		hasAnyPermission, err := models.UserHasAnyPermission(userID, permissionNames, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAnyPermission {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions middleware que verifica si el usuario tiene todos los permisos especificados
func RequireAllPermissions(permissionNames []string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene todos los permisos
		hasAllPermissions, err := models.UserHasAllPermissions(userID, permissionNames, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasAllPermissions {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckRoleOrPermission middleware que verifica si el usuario tiene un rol O un permiso específico
func CheckRoleOrPermission(roleName string, permissionName string, guardName ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si el usuario está autenticado
		if !utils.IsUserAuthenticated(c.Request) {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You must be logged in to access this page.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Obtener el usuario autenticado de la sesión
		userID, authenticated := getUserFromSession(c)
		if !authenticated {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Session expired. Please log in again.")
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		// Determinar el guard name
		guard := "web"
		if len(guardName) > 0 && guardName[0] != "" {
			guard = guardName[0]
		}

		// Verificar si el usuario tiene el rol o el permiso
		hasRole, err := models.UserHasRoleByName(userID, roleName, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if hasRole {
			c.Next()
			return
		}

		hasPermission, err := models.UserHasPermission(userID, permissionName, guard)
		if err != nil {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error checking user permissions.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.CreateFlashNotification(c.Writer, c.Request, "error", "You don't have permission to access this page.")
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
