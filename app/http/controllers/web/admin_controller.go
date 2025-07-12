package web

import (
	"net/http"
	"strconv"
	"web_utilidades/app/helpers"
	"web_utilidades/app/models"
	"web_utilidades/app/structs"
	"web_utilidades/app/utils"

	"github.com/gin-gonic/gin"
)

// AdminController maneja las operaciones administrativas del panel
type AdminController struct{}

// Dashboard muestra el panel de administración
func (ac *AdminController) Dashboard(c *gin.Context) {
	// Verificar si el usuario puede acceder al dashboard
	if !helpers.CanAccessDashboard(c.Request) {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "No tienes permisos para acceder al dashboard.")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Obtener información del usuario autenticado
	user, _ := utils.GetAuthenticatedUser(c.Request)

	// Preparar datos para la vista
	data := gin.H{
		"user":             user,
		"can_manage_users": helpers.CanManageUsers(c.Request),
		"can_manage_roles": helpers.CanManageRoles(c.Request),
		"can_manage_perms": helpers.CanManagePermissions(c.Request),
		"is_admin":         helpers.IsUserAdmin(c.Request),
		"is_super_admin":   helpers.IsUserSuperAdmin(c.Request),
	}

	// Obtener roles y permisos del usuario para mostrar en la vista
	userRoles, _ := helpers.GetUserRoles(c.Request)
	userPermissions, _ := helpers.GetUserPermissions(c.Request)
	data["user_roles"] = userRoles
	data["user_permissions"] = userPermissions

	c.HTML(http.StatusOK, "admin/dashboard.html", data)
}

// UsersIndex muestra la lista de usuarios (solo para quienes tengan permiso)
func (ac *AdminController) UsersIndex(c *gin.Context) {
	// Verificar permiso para ver usuarios
	if !helpers.HasPermission(c.Request, "view-users") {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "No tienes permisos para ver usuarios.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener usuarios: "+err.Error())
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	data := gin.H{
		"users":            users,
		"can_create_users": helpers.HasPermission(c.Request, "create-users"),
		"can_edit_users":   helpers.HasPermission(c.Request, "edit-users"),
		"can_delete_users": helpers.HasPermission(c.Request, "delete-users"),
		"can_assign_roles": helpers.HasPermission(c.Request, "assign-roles"),
	}

	c.HTML(http.StatusOK, "admin/users/index.html", data)
}

// UserShow muestra un usuario específico con sus roles y permisos
func (ac *AdminController) UserShow(c *gin.Context) {
	// Verificar permiso para ver usuarios
	if !helpers.HasPermission(c.Request, "view-users") {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "No tienes permisos para ver usuarios.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	userID := c.Param("id")
	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Usuario no encontrado.")
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	userIDInt, _ := strconv.Atoi(userID)

	// Obtener roles y permisos del usuario
	userRoles, err := models.GetUserRoles(userIDInt)
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener roles del usuario.")
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	userDirectPermissions, err := models.GetUserDirectPermissions(userIDInt)
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener permisos del usuario.")
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	userAllPermissions, err := models.GetUserAllPermissions(userIDInt)
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener todos los permisos del usuario.")
		c.Redirect(http.StatusSeeOther, "/admin/users")
		return
	}

	// Obtener todos los roles disponibles para asignación
	availableRoles, err := models.GetAllRoles()
	if err != nil {
		availableRoles = []structs.RoleStruct{}
	}

	// Obtener todos los permisos disponibles para asignación directa
	availablePermissions, err := models.GetAllPermissions()
	if err != nil {
		availablePermissions = []structs.PermissionStruct{}
	}

	data := gin.H{
		"user":                    user,
		"user_roles":              userRoles,
		"user_direct_permissions": userDirectPermissions,
		"user_all_permissions":    userAllPermissions,
		"available_roles":         availableRoles,
		"available_permissions":   availablePermissions,
		"can_edit_users":          helpers.HasPermission(c.Request, "edit-users"),
		"can_assign_roles":        helpers.HasPermission(c.Request, "assign-roles"),
		"can_assign_permissions":  helpers.HasPermission(c.Request, "assign-permissions"),
	}

	c.HTML(http.StatusOK, "admin/users/show.html", data)
}

// RolesIndex muestra la lista de roles
func (ac *AdminController) RolesIndex(c *gin.Context) {
	// Verificar permiso para ver roles
	if !helpers.HasPermission(c.Request, "view-roles") {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "No tienes permisos para ver roles.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	roles, err := models.GetAllRoles()
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener roles: "+err.Error())
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	data := gin.H{
		"roles":            roles,
		"can_create_roles": helpers.HasPermission(c.Request, "create-roles"),
		"can_edit_roles":   helpers.HasPermission(c.Request, "edit-roles"),
		"can_delete_roles": helpers.HasPermission(c.Request, "delete-roles"),
		"can_assign_roles": helpers.HasPermission(c.Request, "assign-roles"),
	}

	c.HTML(http.StatusOK, "admin/roles/index.html", data)
}

// PermissionsIndex muestra la lista de permisos
func (ac *AdminController) PermissionsIndex(c *gin.Context) {
	// Verificar permiso para ver permisos
	if !helpers.HasPermission(c.Request, "view-permissions") {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "No tienes permisos para ver permisos.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	permissions, err := models.GetAllPermissions()
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error al obtener permisos: "+err.Error())
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	data := gin.H{
		"permissions":            permissions,
		"can_create_permissions": helpers.HasPermission(c.Request, "create-permissions"),
		"can_edit_permissions":   helpers.HasPermission(c.Request, "edit-permissions"),
		"can_delete_permissions": helpers.HasPermission(c.Request, "delete-permissions"),
		"can_assign_permissions": helpers.HasPermission(c.Request, "assign-permissions"),
	}

	c.HTML(http.StatusOK, "admin/permissions/index.html", data)
}

// Ejemplo de uso con múltiples verificaciones
func (ac *AdminController) AdvancedPermissionExample(c *gin.Context) {
	// Verificar múltiples condiciones
	if !helpers.IsUserAdmin(c.Request) {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Solo administradores pueden acceder.")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Verificar rol específico O permiso específico
	if !helpers.HasRole(c.Request, "super-admin") && !helpers.HasPermission(c.Request, "manage-settings") {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Necesitas ser super-admin o tener el permiso 'manage-settings'.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	// Verificar cualquiera de varios roles
	requiredRoles := []string{"admin", "super-admin", "editor"}
	if !helpers.HasAnyRole(c.Request, requiredRoles) {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Necesitas al menos uno de estos roles: admin, super-admin, editor.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	// Verificar todos los permisos requeridos
	requiredPermissions := []string{"view-dashboard", "manage-settings"}
	if !helpers.HasAllPermissions(c.Request, requiredPermissions) {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Necesitas todos estos permisos: view-dashboard, manage-settings.")
		c.Redirect(http.StatusSeeOther, "/admin")
		return
	}

	c.HTML(http.StatusOK, "admin/advanced.html", gin.H{
		"message": "¡Tienes todos los permisos necesarios!",
	})
}
