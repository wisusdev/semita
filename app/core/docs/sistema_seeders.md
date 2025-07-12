# Cambios en el Sistema de Seeders

## Resumen de Cambios

Se ha modificado el comportamiento del sistema de seeders para que **siempre limpie y vuelva a poblar las tablas** en lugar de verificar si ya fueron ejecutados previamente.

## Comportamiento Anterior vs Nuevo

### Anterior
- Los seeders verificaban en la tabla `seeders` si ya habían sido ejecutados
- Si ya estaban ejecutados, se saltaban la ejecución
- Era necesario hacer rollback manual antes de volver a ejecutar

### Nuevo
- Los seeders **siempre** se ejecutan limpiando primero los datos existentes
- No se verifica el estado previo de ejecución
- Cada ejecución garantiza datos frescos y consistentes
- La tabla `seeders` se mantiene solo para auditoría/logs (opcional)

## Comandos Disponibles

### `go run main.go seed all`
Ejecuta todos los seeders registrados:
1. Limpia automáticamente los datos existentes (rollback)
2. Ejecuta el seeding con datos frescos
3. Respeta las dependencias entre seeders

### `go run main.go seed run [nombre_seeder]`
Ejecuta un seeder específico:
1. Ejecuta primero las dependencias requeridas
2. Limpia los datos existentes del seeder
3. Ejecuta el seeding con datos frescos

### `go run main.go seed clean`
**Nuevo comando**: Limpia todos los datos creados por seeders:
- Ejecuta rollback en todos los seeders registrados
- Útil para limpiar completamente la base de datos antes de poblar

### `go run main.go seed rollback [nombre_seeder]`
Ejecuta solo el rollback de un seeder específico

### `go run main.go seed reset [nombre_seeder]`
Equivalente a `run` con el nuevo comportamiento (mantiene compatibilidad)

### `go run main.go seed status`
Muestra el estado de ejecución de todos los seeders (solo para auditoría)

## Cambios Técnicos Implementados

### 1. SeederManager.IsSeederExecuted()
```go
// Antes: Verificaba la tabla seeders
// Ahora: Siempre retorna false para forzar re-ejecución
func (sm *SeederManager) IsSeederExecuted(name string) bool {
    return false // Siempre ejecutar
}
```

### 2. SeederManager.RunSeeder()
```go
// Nuevo comportamiento:
// 1. Ejecuta dependencias
// 2. Hace rollback para limpiar datos
// 3. Ejecuta el seeding
// 4. Registra en logs (opcional)
```

### 3. SeederManager.MarkSeederAsExecuted()
```go
// Ahora es solo para auditoría/logs
// No afecta la lógica de ejecución
// Los errores no son fatales
```

### 4. Nuevo método CleanAllSeederData()
```go
// Limpia todos los datos de seeders en orden inverso
// Respeta las dependencias para evitar errores de FK
```

## Flujo de Ejecución

### Seeders Individuales
```
seed run users_seeder
    ↓
1. Ejecutar dependencias (roles_permissions_seeder)
    ↓
2. Limpiar datos existentes (rollback)
    ↓
3. Ejecutar seeding con datos frescos
    ↓
4. Registrar ejecución (opcional para logs)
```

### Todos los Seeders
```
seed all
    ↓
1. Determinar orden por dependencias
    ↓
2. Para cada seeder:
   - Limpiar datos existentes
   - Ejecutar seeding
    ↓
3. Datos completamente frescos en todas las tablas
```

## Ventajas del Nuevo Comportamiento

1. **Consistencia garantizada**: Cada ejecución produce el mismo resultado
2. **Datos frescos**: No hay datos obsoletos o duplicados
3. **Simplicidad**: No hay que gestionar manualmente el estado de seeders
4. **Desarrollo ágil**: Ideal para desarrollo donde se cambian datos frecuentemente
5. **Testing**: Perfecto para tests que requieren datos limpios

## Consideraciones

1. **Tiempo de ejecución**: Puede ser ligeramente más lento ya que siempre limpia datos
2. **Datos de producción**: Usar con cuidado en producción, siempre limpia datos existentes
3. **Transacciones**: Todo se ejecuta en transacciones para garantizar consistencia
4. **Logs**: La tabla `seeders` se mantiene para auditoría pero no afecta la ejecución

## Migración

No se requieren cambios en los seeders existentes. El comportamiento es compatible hacia atrás, solo cambia la lógica interna de ejecución.

Los comandos existentes mantienen la misma interfaz pero con el nuevo comportamiento de limpieza automática.
