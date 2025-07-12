package seeders

import (
	"log"
	"web_utilidades/app/core/database"
	"web_utilidades/app/models"
	"web_utilidades/app/structs"
	"web_utilidades/config"
)

// RolesPermissionsSeeder seeder para roles y permisos
type RolesPermissionsSeeder struct {
	database.BaseSeeder
}

// NewRolesPermissionsSeeder crea una nueva instancia del seeder
func NewRolesPermissionsSeeder() *RolesPermissionsSeeder {
	return &RolesPermissionsSeeder{
		BaseSeeder: database.BaseSeeder{
			DB:   config.DatabaseConnect(),
			Name: "roles_permissions_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (rps *RolesPermissionsSeeder) GetName() string {
	return rps.Name
}

// GetDependencies retorna las dependencias del seeder
func (rps *RolesPermissionsSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// Seed ejecuta el seeding de roles y permisos
func (rps *RolesPermissionsSeeder) Seed() error {
	log.Println("Seeding roles and permissions...")

	// Crear permisos básicos
	permissions := []structs.CreatePermissionStruct{
		{Name: "create-users", GuardName: "web", Description: "Crear usuarios"},
		{Name: "edit-users", GuardName: "web", Description: "Editar usuarios"},
		{Name: "delete-users", GuardName: "web", Description: "Eliminar usuarios"},
		{Name: "view-users", GuardName: "web", Description: "Ver usuarios"},
		{Name: "create-roles", GuardName: "web", Description: "Crear roles"},
		{Name: "edit-roles", GuardName: "web", Description: "Editar roles"},
		{Name: "delete-roles", GuardName: "web", Description: "Eliminar roles"},
		{Name: "view-roles", GuardName: "web", Description: "Ver roles"},
		{Name: "assign-roles", GuardName: "web", Description: "Asignar roles"},
		{Name: "create-permissions", GuardName: "web", Description: "Crear permisos"},
		{Name: "edit-permissions", GuardName: "web", Description: "Editar permisos"},
		{Name: "delete-permissions", GuardName: "web", Description: "Eliminar permisos"},
		{Name: "view-permissions", GuardName: "web", Description: "Ver permisos"},
		{Name: "assign-permissions", GuardName: "web", Description: "Asignar permisos"},
		{Name: "manage-posts", GuardName: "web", Description: "Gestionar posts"},
		{Name: "publish-posts", GuardName: "web", Description: "Publicar posts"},
		{Name: "edit-posts", GuardName: "web", Description: "Editar posts"},
		{Name: "delete-posts", GuardName: "web", Description: "Eliminar posts"},
		{Name: "view-dashboard", GuardName: "web", Description: "Ver dashboard administrativo"},
		{Name: "manage-settings", GuardName: "web", Description: "Gestionar configuración del sistema"},
	}

	log.Println("Creating permissions...")
	createdPermissions := make(map[string]*structs.PermissionStruct)
	for _, permData := range permissions {
		// Verificar si el permiso ya existe
		existingPerm, err := models.GetPermissionByName(permData.Name, permData.GuardName)
		if err == nil {
			log.Printf("Permission '%s' already exists, skipping...", permData.Name)
			createdPermissions[permData.Name] = existingPerm
			continue
		}

		permission, err := models.CreatePermission(permData)
		if err != nil {
			log.Printf("Error creating permission '%s': %v", permData.Name, err)
			continue
		}
		createdPermissions[permData.Name] = permission
		log.Printf("Created permission: %s", permission.Name)
	}

	// Crear roles básicos
	roles := []structs.CreateRoleStruct{
		{Name: "super-admin", GuardName: "web", Description: "Super administrador con todos los permisos"},
		{Name: "admin", GuardName: "web", Description: "Administrador del sistema"},
		{Name: "editor", GuardName: "web", Description: "Editor de contenido"},
		{Name: "moderator", GuardName: "web", Description: "Moderador"},
		{Name: "user", GuardName: "web", Description: "Usuario regular"},
	}

	log.Println("Creating roles...")
	createdRoles := make(map[string]*structs.RoleStruct)
	for _, roleData := range roles {
		// Verificar si el rol ya existe
		existingRole, err := models.GetRoleByName(roleData.Name, roleData.GuardName)
		if err == nil {
			log.Printf("Role '%s' already exists, skipping...", roleData.Name)
			createdRoles[roleData.Name] = existingRole
			continue
		}

		role, err := models.CreateRole(roleData)
		if err != nil {
			log.Printf("Error creating role '%s': %v", roleData.Name, err)
			continue
		}
		createdRoles[roleData.Name] = role
		log.Printf("Created role: %s", role.Name)
	}

	// Asignar permisos a roles
	log.Println("Assigning permissions to roles...")

	// Super Admin - todos los permisos
	if superAdmin, exists := createdRoles["super-admin"]; exists {
		for _, permission := range createdPermissions {
			err := models.AssignPermissionToRole(superAdmin.ID, permission.ID)
			if err != nil {
				log.Printf("Error assigning permission '%s' to role 'super-admin': %v", permission.Name, err)
			}
		}
		log.Println("Assigned all permissions to super-admin role")
	}

	// Admin - permisos administrativos
	if admin, exists := createdRoles["admin"]; exists {
		adminPermissions := []string{
			"create-users", "edit-users", "delete-users", "view-users",
			"create-roles", "edit-roles", "view-roles", "assign-roles",
			"view-permissions", "assign-permissions",
			"manage-posts", "publish-posts", "edit-posts", "delete-posts",
			"view-dashboard", "manage-settings",
		}
		for _, permName := range adminPermissions {
			if permission, exists := createdPermissions[permName]; exists {
				err := models.AssignPermissionToRole(admin.ID, permission.ID)
				if err != nil {
					log.Printf("Error assigning permission '%s' to role 'admin': %v", permission.Name, err)
				}
			}
		}
		log.Println("Assigned admin permissions to admin role")
	}

	// Editor - permisos de contenido
	if editor, exists := createdRoles["editor"]; exists {
		editorPermissions := []string{
			"view-users",
			"manage-posts", "publish-posts", "edit-posts",
			"view-dashboard",
		}
		for _, permName := range editorPermissions {
			if permission, exists := createdPermissions[permName]; exists {
				err := models.AssignPermissionToRole(editor.ID, permission.ID)
				if err != nil {
					log.Printf("Error assigning permission '%s' to role 'editor': %v", permission.Name, err)
				}
			}
		}
		log.Println("Assigned editor permissions to editor role")
	}

	// Moderator - permisos básicos de moderación
	if moderator, exists := createdRoles["moderator"]; exists {
		moderatorPermissions := []string{
			"view-users",
			"edit-posts",
			"view-dashboard",
		}
		for _, permName := range moderatorPermissions {
			if permission, exists := createdPermissions[permName]; exists {
				err := models.AssignPermissionToRole(moderator.ID, permission.ID)
				if err != nil {
					log.Printf("Error assigning permission '%s' to role 'moderator': %v", permission.Name, err)
				}
			}
		}
		log.Println("Assigned moderator permissions to moderator role")
	}

	log.Println("Roles and permissions seeding completed successfully!")
	return nil
}

// Rollback revierte el seeding de roles y permisos
func (rps *RolesPermissionsSeeder) Rollback() error {
	log.Println("Rolling back roles and permissions...")

	// Eliminar todas las relaciones role_permissions
	query := `DELETE FROM role_permissions`
	_, err := rps.DB.Exec(query)
	if err != nil {
		log.Printf("Error deleting role permissions: %v", err)
		return err
	}

	// Eliminar todas las relaciones user_roles
	query = `DELETE FROM user_roles`
	_, err = rps.DB.Exec(query)
	if err != nil {
		log.Printf("Error deleting user roles: %v", err)
		return err
	}

	// Eliminar todas las relaciones user_permissions
	query = `DELETE FROM user_permissions`
	_, err = rps.DB.Exec(query)
	if err != nil {
		log.Printf("Error deleting user permissions: %v", err)
		return err
	}

	// Eliminar roles
	roleNames := []string{"super-admin", "admin", "editor", "moderator", "user"}
	for _, roleName := range roleNames {
		query = `DELETE FROM roles WHERE name = ? AND guard_name = 'web'`
		_, err = rps.DB.Exec(query, roleName)
		if err != nil {
			log.Printf("Error deleting role '%s': %v", roleName, err)
		} else {
			log.Printf("Deleted role: %s", roleName)
		}
	}

	// Eliminar permisos
	permissionNames := []string{
		"create-users", "edit-users", "delete-users", "view-users",
		"create-roles", "edit-roles", "delete-roles", "view-roles", "assign-roles",
		"create-permissions", "edit-permissions", "delete-permissions", "view-permissions", "assign-permissions",
		"manage-posts", "publish-posts", "edit-posts", "delete-posts",
		"view-dashboard", "manage-settings",
	}
	for _, permName := range permissionNames {
		query = `DELETE FROM permissions WHERE name = ? AND guard_name = 'web'`
		_, err = rps.DB.Exec(query, permName)
		if err != nil {
			log.Printf("Error deleting permission '%s': %v", permName, err)
		} else {
			log.Printf("Deleted permission: %s", permName)
		}
	}

	log.Println("Roles and permissions rollback completed successfully!")
	return nil
}
