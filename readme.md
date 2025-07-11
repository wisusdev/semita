# Comandos disponibles

Para manejar migraciones de base de datos en Go, ahora puedes usar la consola tipo Artisan:

Comando para ejecutar el servidor web:

```bash
go run .
```

Comando para ejecutar migraciones:

```bash
go run . migrate
```

Revertir último lote de migraciones:

```bash
go run . migrate:rollback
```

Refrescar la base de datos:

```bash
go run . migrate:fresh
```

## Comandos Artisan disponibles

- Generar clave JWT:

```bash
go run . key:generate
```

- Generar llaves OAuth2 (RSA):

```bash
go run . oauth:keys
```

- Crear una nueva migración:

```bash
go run . make:migration NombreDeLaMigracion
```

Esto creará un archivo en `database/migrations` con el formato `2024_06_10_123456_NombreDeLaMigracion.go`.

Luego, edita el archivo generado para definir la lógica de creación y reversión de la tabla.

## JSON API

Para ejecutar un servidor JSON que sirva como base de datos de prueba:

```bash
npm install -g json-server
json-server --watch ./database/json/db.json
# Accede a http://localhost:3000/
```

Compilar proyecto

```bash
go build
```
