# Comando make:migration-from-db

Este comando genera automáticamente archivos de migración en Go a partir de una base de datos MySQL existente.

## Uso

```bash
go run main.go make:migration-from-db
```

## Funcionalidades

- **Análisis automático**: Conecta a la base de datos configurada en `.env` y analiza todas las tablas existentes
- **Generación de migraciones**: Crea archivos de migración que siguen exactamente la estructura y convenciones del proyecto
- **Detección de dependencias**: Ordena las tablas según sus relaciones de claves foráneas
- **Preservación de estructura**: Mantiene todos los tipos de datos, constraints, índices y valores por defecto

## Características generadas

### Estructura de archivos

- Nombre: `YYYY_MM_DD_HHMMSS_create_[tabla]_table.go`
- Ubicación: `database/migrations/`
- Package: `migrations`

### Estructura de clases

- Struct: `Create[Tabla]Table` que embebe `database.BaseMigration`
- Constructor: `NewCreate[Tabla]Table()`
- Métodos: `Up(db *sql.DB) error` y `Down(db *sql.DB) error`

### Elementos detectados y generados

- **Columnas**: Todos los tipos de datos MySQL con sus constraints
- **Claves primarias**: AUTO_INCREMENT, PRIMARY KEY
- **Índices**: UNIQUE, INDEX con múltiples columnas
- **Claves foráneas**: FOREIGN KEY con ON DELETE/UPDATE rules
- **Valores por defecto**: DEFAULT values, CURRENT_TIMESTAMP
- **Constraints**: NOT NULL, UNIQUE

## Configuración

El comando usa la configuración de base de datos del archivo `.env`:

```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=tu_base_datos
DB_USER=tu_usuario
DB_PASSWORD=tu_password
```

## Exclusiones

- La tabla `migrations` es automáticamente excluida del proceso de generación

## Ejemplo de salida

Para una tabla `users`, genera:

```go
package migrations

import (
	"database/sql"
	"semita/app/core/database"
)

type CreateUsersTable struct {
	database.BaseMigration
}

func NewCreateUsersTable() *CreateUsersTable {
	return &CreateUsersTable{
		BaseMigration: database.BaseMigration{
			Name:      "create_users_table",
			Timestamp: "2025_07_13_000001",
		},
	}
}

func (m *CreateUsersTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			email_verified_at DATETIME,
			remember_token VARCHAR(100),
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func (m *CreateUsersTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}
```

## Notas importantes

- Los archivos se generan con timestamps únicos incrementales
- Las tablas se ordenan automáticamente según dependencias de claves foráneas
- Los archivos existentes no se sobrescriben
- Se mantiene compatibilidad total con el sistema de migraciones existente
