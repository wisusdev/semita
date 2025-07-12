package main

import (
	"fmt"
	"log"
	"web_utilidades/app/models"
	"web_utilidades/app/structs"
	"web_utilidades/config"
)

// Este archivo muestra ejemplos de uso del sistema de roles y permisos
func main() {
	// Conectar a la base de datos
	config.DatabaseConnect()

	// Ejemplo 1: Crear un rol
	fmt.Println("=== Ejemplo 1: Crear rol ===")
	roleData := structs.CreateRoleStruct{
		Name:        "content-manager",
		GuardName:   "web",
		Description: "Gestor de contenido",
	}

	role, err := models.CreateRole(roleData)
	if err != nil {
		log.Printf("Error creando rol: %v", err)
	} else {
		fmt.Printf("Rol creado: %s (ID: %d)\n", role.Name, role.ID)
	}

	// Ejemplo 2: Crear un permiso
	fmt.Println("\n=== Ejemplo 2: Crear permiso ===")
	permissionData := structs.CreatePermissionStruct{
		Name:        "moderate-comments",
		GuardName:   "web",
		Description: "Moderar comentarios",
	}

	permission, err := models.CreatePermission(permissionData)
	if err != nil {
		log.Printf("Error creando permiso: %v", err)
	} else {
		fmt.Printf("Permiso creado: %s (ID: %d)\n", permission.Name, permission.ID)
	}

	// Ejemplo 3: Asignar permiso a rol
	fmt.Println("\n=== Ejemplo 3: Asignar permiso a rol ===")
	if role != nil && permission != nil {
		err = models.AssignPermissionToRole(role.ID, permission.ID)
		if err != nil {
			log.Printf("Error asignando permiso a rol: %v", err)
		} else {
			fmt.Printf("Permiso '%s' asignado al rol '%s'\n", permission.Name, role.Name)
		}
	}

	// Ejemplo 4: Asignar rol a usuario (usuario ID 1)
	fmt.Println("\n=== Ejemplo 4: Asignar rol a usuario ===")
	userID := 1
	if role != nil {
		err = models.AssignRoleToUser(userID, role.ID)
		if err != nil {
			log.Printf("Error asignando rol a usuario: %v", err)
		} else {
			fmt.Printf("Rol '%s' asignado al usuario ID %d\n", role.Name, userID)
		}
	}

	// Ejemplo 5: Verificar si usuario tiene rol
	fmt.Println("\n=== Ejemplo 5: Verificar rol de usuario ===")
	hasRole, err := models.UserHasRoleByName(userID, "content-manager", "web")
	if err != nil {
		log.Printf("Error verificando rol: %v", err)
	} else {
		fmt.Printf("Usuario %d tiene rol 'content-manager': %t\n", userID, hasRole)
	}

	// Ejemplo 6: Verificar si usuario tiene permiso
	fmt.Println("\n=== Ejemplo 6: Verificar permiso de usuario ===")
	hasPermission, err := models.UserHasPermission(userID, "moderate-comments", "web")
	if err != nil {
		log.Printf("Error verificando permiso: %v", err)
	} else {
		fmt.Printf("Usuario %d tiene permiso 'moderate-comments': %t\n", userID, hasPermission)
	}

	// Ejemplo 7: Obtener todos los roles del usuario
	fmt.Println("\n=== Ejemplo 7: Obtener roles del usuario ===")
	userRoles, err := models.GetUserRoles(userID)
	if err != nil {
		log.Printf("Error obteniendo roles: %v", err)
	} else {
		fmt.Printf("Usuario %d tiene %d roles:\n", userID, len(userRoles))
		for _, userRole := range userRoles {
			fmt.Printf("- %s: %s\n", userRole.Name, userRole.Description)
		}
	}

	// Ejemplo 8: Obtener todos los permisos del usuario
	fmt.Println("\n=== Ejemplo 8: Obtener permisos del usuario ===")
	userPermissions, err := models.GetUserAllPermissions(userID)
	if err != nil {
		log.Printf("Error obteniendo permisos: %v", err)
	} else {
		fmt.Printf("Usuario %d tiene %d permisos:\n", userID, len(userPermissions))
		for _, userPerm := range userPermissions {
			fmt.Printf("- %s: %s\n", userPerm.Name, userPerm.Description)
		}
	}

	// Ejemplo 9: Verificar múltiples roles
	fmt.Println("\n=== Ejemplo 9: Verificar múltiples roles ===")
	roleNames := []string{"admin", "content-manager", "editor"}
	hasAnyRole, err := models.UserHasAnyRole(userID, roleNames, "web")
	if err != nil {
		log.Printf("Error verificando múltiples roles: %v", err)
	} else {
		fmt.Printf("Usuario %d tiene al menos uno de estos roles %v: %t\n", userID, roleNames, hasAnyRole)
	}

	// Ejemplo 10: Asignar permiso directo a usuario
	fmt.Println("\n=== Ejemplo 10: Asignar permiso directo a usuario ===")
	// Crear un permiso especial
	specialPermission := structs.CreatePermissionStruct{
		Name:        "special-access",
		GuardName:   "web",
		Description: "Acceso especial",
	}

	specialPerm, err := models.CreatePermission(specialPermission)
	if err != nil {
		log.Printf("Error creando permiso especial: %v", err)
	} else {
		// Asignar directamente al usuario
		err = models.AssignPermissionToUser(userID, specialPerm.ID)
		if err != nil {
			log.Printf("Error asignando permiso directo: %v", err)
		} else {
			fmt.Printf("Permiso '%s' asignado directamente al usuario %d\n", specialPerm.Name, userID)
		}
	}

	// Ejemplo 11: Verificar permisos directos vs todos los permisos
	fmt.Println("\n=== Ejemplo 11: Permisos directos vs todos los permisos ===")
	directPerms, err := models.GetUserDirectPermissions(userID)
	if err != nil {
		log.Printf("Error obteniendo permisos directos: %v", err)
	} else {
		fmt.Printf("Permisos directos (%d):\n", len(directPerms))
		for _, perm := range directPerms {
			fmt.Printf("- %s (directo)\n", perm.Name)
		}
	}

	allPerms, err := models.GetUserAllPermissions(userID)
	if err != nil {
		log.Printf("Error obteniendo todos los permisos: %v", err)
	} else {
		fmt.Printf("Todos los permisos (%d):\n", len(allPerms))
		for _, perm := range allPerms {
			fmt.Printf("- %s\n", perm.Name)
		}
	}

	fmt.Println("\n=== Ejemplos completados ===")
	fmt.Println("Para usar este sistema en tu aplicación:")
	fmt.Println("1. Ejecuta las migraciones: go run main.go migrate")
	fmt.Println("2. Ejecuta el seeder: go run main.go seed:roles-permissions")
	fmt.Println("3. Usa los middleware en tus rutas")
	fmt.Println("4. Usa los helpers en tus controladores")
	fmt.Println("5. Consulta la documentación en docs/roles_permisos_sistema.md")
}
