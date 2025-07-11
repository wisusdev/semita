package models

import (
	"time"
	"web_utilidades/app/structs"
	"web_utilidades/config"
)

var userTable = "users"

func GetAllUsers() ([]structs.UserStruct, error) {
	// Instanciamos la conexión a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener todos los usuarios
	var query = "SELECT id, name, email, created_at, updated_at FROM " + userTable

	// Ejecutamos la consulta y obtenemos los resultados
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Creamos un slice para almacenar los usuarios
	var users []structs.UserStruct

	// Iteramos sobre los resultados y los agregamos al slice
	for rows.Next() {
		var user structs.UserStruct
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func StoreUser(user structs.StoreUserStruct) (err error) {
	// Instanciamos la conexion a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexion se cierre al final de la funcion
	defer database.Close()

	// Preparamos la consulta para insertar un nuevo usuario
	var query = "INSERT INTO " + userTable + " (name, email, password) VALUES (?, ?, ?)"

	// Ejecutamos la consulta con los datos del usuario
	_, err = database.Exec(query, user.Name, user.Email, user.Password)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id string) (user structs.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su ID
	var query = "SELECT id, name, email, password, created_at, updated_at FROM " + userTable + " WHERE id = ?"

	// Ejecutamos la consulta y obtenemos los resultados
	err = database.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// Si hubo un error al ejecutar la consulta o no se encontró el usuario, retornamos el error
	if err != nil {
		return structs.UserStruct{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (user structs.UserStruct, err error) {
	// Instanciamos la conexión a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para obtener un usuario por su email
	var query = "SELECT id, name, email, password, created_at, updated_at FROM " + userTable + " WHERE email = ?"

	// Ejecutamos la consulta y obtenemos los resultados
	err = database.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// Si hubo un error al ejecutar la consulta o no se encontró el usuario, retornamos el error
	if err != nil {
		return structs.UserStruct{}, err
	}

	return user, nil
}

func UpdateUser(user structs.UpdateUserStruct) (err error) {
	// Instanciamos la conexión a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para actualizar un usuario por su ID
	var query = "UPDATE " + userTable + " SET name = ?, email = ?, password = ? WHERE id = ?"

	// Ejecutamos la consulta con los datos del usuario
	_, err = database.Exec(query, user.Name, user.Email, user.Password, user.ID)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string) (err error) {
	// Instanciamos la conexión a la base de datos
	var database = config.DatabaseConnect()

	// Aseguramos que la conexión se cierre al final de la función
	defer database.Close()

	// Preparamos la consulta para eliminar un usuario por su ID
	var query = "DELETE FROM " + userTable + " WHERE id = ?"

	// Ejecutamos la consulta con el ID del usuario
	_, err = database.Exec(query, id)

	// Si hubo un error al ejecutar la consulta, retornamos el error
	if err != nil {
		return err
	}

	return nil
}

// MarkEmailVerified actualiza el campo email_verified_at del usuario
func MarkEmailVerified(userID int) error {
	db := config.DatabaseConnect()
	defer db.Close()
	_, err := db.Exec("UPDATE "+userTable+" SET email_verified_at = ? WHERE id = ?", time.Now().Format("2006-01-02 15:04:05"), userID)
	return err
}
