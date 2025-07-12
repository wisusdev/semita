package base

import (
	"net/http"
	"strconv"
	"web_utilidades/app/models"
	"web_utilidades/app/structs"
	"web_utilidades/app/utils"

	"github.com/gin-gonic/gin"
)

// RoleController maneja las operaciones CRUD de roles
type RoleController struct{}

// Index muestra todos los roles
func (rc *RoleController) Index(c *gin.Context) {
	roles, err := models.GetAllRoles()
	if err != nil {
		utils.CreateFlashNotification(c.Writer, c.Request, "error", "Error retrieving roles: "+err.Error())
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}

// Show muestra un rol específico con sus permisos
func (rc *RoleController) Show(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	role, err := models.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Role not found",
		})
		return
	}

	permissions, err := models.GetRolePermissions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving role permissions: " + err.Error(),
		})
		return
	}

	roleWithPermissions := structs.RoleWithPermissions{
		RoleStruct:  *role,
		Permissions: permissions,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roleWithPermissions,
	})
}

// Store crea un nuevo rol
func (rc *RoleController) Store(c *gin.Context) {
	var roleData structs.CreateRoleStruct
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	role, err := models.CreateRole(roleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Role created successfully",
		"data":    role,
	})
}

// Update actualiza un rol existente
func (rc *RoleController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	var roleData structs.CreateRoleStruct
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	role, err := models.UpdateRole(id, roleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error updating role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role updated successfully",
		"data":    role,
	})
}

// Delete elimina un rol
func (rc *RoleController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid role ID",
		})
		return
	}

	err = models.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error deleting role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role deleted successfully",
	})
}

// AssignToUser asigna un rol a un usuario
func (rc *RoleController) AssignToUser(c *gin.Context) {
	var request structs.AssignRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	err := models.AssignRoleToUser(request.UserID, request.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error assigning role to user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role assigned to user successfully",
	})
}

// RevokeFromUser revoca un rol de un usuario
func (rc *RoleController) RevokeFromUser(c *gin.Context) {
	var request structs.AssignRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"errors":  err.Error(),
		})
		return
	}

	err := models.RevokeRoleFromUser(request.UserID, request.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error revoking role from user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Role revoked from user successfully",
	})
}

// GetUserRoles obtiene todos los roles de un usuario
func (rc *RoleController) GetUserRoles(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	roles, err := models.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error retrieving user roles: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}
