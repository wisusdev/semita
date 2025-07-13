package seeders

import (
	"database/sql"
	"log"
	"web_utilidades/app/core/database"
	"web_utilidades/config"
)

// UsersSeeder seeder para usuarios de prueba
type UsersSeeder struct {
	database.BaseSeeder
}

// NewUsersSeeder crea una nueva instancia del seeder
func NewUsersSeeder() *UsersSeeder {
	return &UsersSeeder{
		BaseSeeder: database.BaseSeeder{
			DB:   config.DatabaseConnect(),
			Name: "users_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (us *UsersSeeder) GetName() string {
	return us.Name
}

// GetDependencies retorna las dependencias del seeder
func (us *UsersSeeder) GetDependencies() []string {
	return []string{"roles_permissions_seeder"} // Depende de roles y permisos
}

// Seed ejecuta el seeding de usuarios
func (us *UsersSeeder) Seed() error {
	log.Println("Seeding users...")

	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{
			Name:     "Super Administrador",
			Email:    "superadmin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "super-admin",
		},
		{
			Name:     "Administrador",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "admin",
		},
		{
			Name:     "Editor Principal",
			Email:    "editor@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "editor",
		},
		{
			Name:     "Moderador",
			Email:    "moderator@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "moderator",
		},
		{
			Name:     "Usuario Demo",
			Email:    "user@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "user",
		},
		{
			Name:     "María García",
			Email:    "maria.garcia@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "user",
		},
		{
			Name:     "Carlos López",
			Email:    "carlos.lopez@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "user",
		},
		{
			Name:     "Ana Martínez",
			Email:    "ana.martinez@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "editor",
		},
		{
			Name:     "Luis Pérez",
			Email:    "user00@wisus.dev",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "user",
		},
	}

	for _, user := range users {
		// Verificar si el usuario ya existe
		var existingID int
		checkQuery := `SELECT id FROM users WHERE email = ?`
		err := us.DB.QueryRow(checkQuery, user.Email).Scan(&existingID)

		if err == sql.ErrNoRows {
			// No existe, crear nuevo usuario
			insertQuery := `
				INSERT INTO users (name, email, password, email_verified_at, created_at, updated_at) 
				VALUES (?, ?, ?, NOW(), NOW(), NOW())`

			result, err := us.DB.Exec(insertQuery, user.Name, user.Email, user.Password)
			if err != nil {
				log.Printf("Error creating user '%s': %v", user.Name, err)
				continue
			}

			userID, _ := result.LastInsertId()
			log.Printf("Created user: %s (ID: %d)", user.Name, userID)

			// Asignar rol al usuario
			err = us.assignRoleToUser(int(userID), user.Role)
			if err != nil {
				log.Printf("Error assigning role '%s' to user '%s': %v", user.Role, user.Name, err)
			} else {
				log.Printf("Assigned role '%s' to user '%s'", user.Role, user.Name)
			}
		} else if err != nil {
			log.Printf("Error checking existing user '%s': %v", user.Name, err)
			continue
		} else {
			log.Printf("User '%s' already exists, skipping...", user.Name)
		}
	}

	log.Println("Users seeding completed successfully!")
	return nil
}

// assignRoleToUser asigna un rol a un usuario
func (us *UsersSeeder) assignRoleToUser(userID int, roleName string) error {
	// Obtener el ID del rol
	var roleID int
	roleQuery := `SELECT id FROM roles WHERE name = ? AND guard_name = 'web'`
	err := us.DB.QueryRow(roleQuery, roleName).Scan(&roleID)
	if err != nil {
		return err
	}

	// Verificar si la relación ya existe
	var existingID int
	checkQuery := `SELECT id FROM user_roles WHERE user_id = ? AND role_id = ?`
	err = us.DB.QueryRow(checkQuery, userID, roleID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// No existe, crear la relación
		insertQuery := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
		_, err = us.DB.Exec(insertQuery, userID, roleID)
		return err
	}

	return err
}

// Rollback revierte el seeding de usuarios
func (us *UsersSeeder) Rollback() error {
	log.Println("Rolling back users...")

	userEmails := []string{
		"superadmin@example.com",
		"admin@example.com",
		"editor@example.com",
		"moderator@example.com",
		"user@example.com",
		"maria.garcia@example.com",
		"carlos.lopez@example.com",
		"ana.martinez@example.com",
	}

	for _, email := range userEmails {
		// Primero obtener el ID del usuario
		var userID int
		userQuery := `SELECT id FROM users WHERE email = ?`
		err := us.DB.QueryRow(userQuery, email).Scan(&userID)

		if err == sql.ErrNoRows {
			continue // Usuario no existe, continuar
		} else if err != nil {
			log.Printf("Error getting user ID for '%s': %v", email, err)
			continue
		}

		// Eliminar relaciones user_roles
		deleteRolesQuery := `DELETE FROM user_roles WHERE user_id = ?`
		_, err = us.DB.Exec(deleteRolesQuery, userID)
		if err != nil {
			log.Printf("Error deleting user roles for '%s': %v", email, err)
		}

		// Eliminar relaciones user_permissions
		deletePermissionsQuery := `DELETE FROM user_permissions WHERE user_id = ?`
		_, err = us.DB.Exec(deletePermissionsQuery, userID)
		if err != nil {
			log.Printf("Error deleting user permissions for '%s': %v", email, err)
		}

		// Eliminar el usuario
		deleteUserQuery := `DELETE FROM users WHERE email = ?`
		result, err := us.DB.Exec(deleteUserQuery, email)
		if err != nil {
			log.Printf("Error deleting user '%s': %v", email, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("Deleted user: %s", email)
		}
	}

	log.Println("Users rollback completed successfully!")
	return nil
}
