package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"web_utilidades/config"

	"github.com/spf13/cobra"
)

// Estructura para representar una tabla
type TableInfo struct {
	Name        string
	Columns     []ColumnInfo
	Indexes     []IndexInfo
	ForeignKeys []ForeignKeyInfo
}

// Estructura para representar una columna
type ColumnInfo struct {
	Name    string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

// Estructura para representar un índice
type IndexInfo struct {
	Name    string
	Unique  bool
	Columns []string
}

// Estructura para representar una clave foránea
type ForeignKeyInfo struct {
	Name             string
	Column           string
	ReferencedTable  string
	ReferencedColumn string
	OnDelete         string
	OnUpdate         string
}

var MakeMigrationFromDbCmd = &cobra.Command{
	Use:   "make:migration-from-db",
	Short: "Genera archivos de migración a partir de una base de datos existente",
	Run: func(cmd *cobra.Command, args []string) {
		generateMigrationsFromDatabase()
	},
}

func generateMigrationsFromDatabase() {
	db := config.DatabaseConnect()
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Obtener todas las tablas
	tables, err := getTables(db)
	if err != nil {
		log.Fatal("Error getting tables:", err)
	}

	// Filtrar la tabla migrations
	var filteredTables []string
	for _, table := range tables {
		if table != "migrations" {
			filteredTables = append(filteredTables, table)
		}
	}

	// Ordenar tablas por dependencias (las que no tienen FK primero)
	orderedTables, err := orderTablesByDependencies(db, filteredTables)
	if err != nil {
		log.Fatal("Error ordering tables:", err)
	}

	// Generar migración para cada tabla
	baseTime := time.Now()
	for i, tableName := range orderedTables {
		// Incrementar el timestamp para cada tabla
		timestamp := baseTime.Add(time.Duration(i) * time.Minute)

		tableInfo, err := getTableInfo(db, tableName)
		if err != nil {
			log.Printf("Error getting info for table %s: %v", tableName, err)
			continue
		}

		err = generateMigrationFile(tableInfo, timestamp)
		if err != nil {
			log.Printf("Error generating migration for table %s: %v", tableName, err)
			continue
		}

		fmt.Printf("Generated migration for table: %s\n", tableName)
	}

	fmt.Println("Migration generation completed!")
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func getTableInfo(db *sql.DB, tableName string) (*TableInfo, error) {
	tableInfo := &TableInfo{Name: tableName}

	// Obtener columnas
	columns, err := getTableColumns(db, tableName)
	if err != nil {
		return nil, err
	}
	tableInfo.Columns = columns

	// Obtener índices
	indexes, err := getTableIndexes(db, tableName)
	if err != nil {
		return nil, err
	}
	tableInfo.Indexes = indexes

	// Obtener claves foráneas
	foreignKeys, err := getTableForeignKeys(db, tableName)
	if err != nil {
		return nil, err
	}
	tableInfo.ForeignKeys = foreignKeys

	return tableInfo, nil
}

func getTableColumns(db *sql.DB, tableName string) ([]ColumnInfo, error) {
	query := fmt.Sprintf("DESCRIBE %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(&col.Name, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func getTableIndexes(db *sql.DB, tableName string) ([]IndexInfo, error) {
	query := fmt.Sprintf("SHOW INDEX FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexMap := make(map[string]*IndexInfo)

	for rows.Next() {
		var table, nonUnique, keyName, seqInIndex, columnName, collation, cardinality, subPart, packed, nullVal, indexType, comment, indexComment, visible, expression sql.NullString

		err := rows.Scan(&table, &nonUnique, &keyName, &seqInIndex, &columnName, &collation, &cardinality, &subPart, &packed, &nullVal, &indexType, &comment, &indexComment, &visible, &expression)
		if err != nil {
			return nil, err
		}

		if keyName.String == "PRIMARY" {
			continue // Skip primary key, handled in column definition
		}

		if _, exists := indexMap[keyName.String]; !exists {
			indexMap[keyName.String] = &IndexInfo{
				Name:    keyName.String,
				Unique:  nonUnique.String == "0",
				Columns: []string{},
			}
		}
		indexMap[keyName.String].Columns = append(indexMap[keyName.String].Columns, columnName.String)
	}

	var indexes []IndexInfo
	for _, idx := range indexMap {
		indexes = append(indexes, *idx)
	}

	return indexes, nil
}

func getTableForeignKeys(db *sql.DB, tableName string) ([]ForeignKeyInfo, error) {
	// Usar SHOW CREATE TABLE para obtener las claves foráneas
	query := fmt.Sprintf("SHOW CREATE TABLE %s", tableName)

	var tableName2, createTable string
	err := db.QueryRow(query).Scan(&tableName2, &createTable)
	if err != nil {
		return nil, fmt.Errorf("error getting CREATE TABLE for %s: %v", tableName, err)
	}

	// Parsear el CREATE TABLE para extraer las claves foráneas
	var foreignKeys []ForeignKeyInfo

	// Buscar patrones como: CONSTRAINT `fk_name` FOREIGN KEY (`column`) REFERENCES `table` (`column`)
	lines := strings.Split(createTable, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "CONSTRAINT") && strings.Contains(line, "FOREIGN KEY") {
			// Parsear la línea de foreign key
			fk := parseForeignKeyLine(line)
			if fk != nil {
				foreignKeys = append(foreignKeys, *fk)
			}
		}
	}

	return foreignKeys, nil
}

func parseForeignKeyLine(line string) *ForeignKeyInfo {
	// Ejemplo: CONSTRAINT `posts_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)

	// Buscar el nombre del constraint
	constraintStart := strings.Index(line, "CONSTRAINT")
	if constraintStart == -1 {
		return nil
	}

	// Extraer componentes usando regex básico
	fk := &ForeignKeyInfo{
		OnDelete: "RESTRICT",
		OnUpdate: "RESTRICT",
	}

	// Buscar CONSTRAINT name
	if idx := strings.Index(line, "CONSTRAINT `"); idx != -1 {
		start := idx + 12
		if end := strings.Index(line[start:], "`"); end != -1 {
			fk.Name = line[start : start+end]
		}
	}

	// Buscar FOREIGN KEY (column)
	if idx := strings.Index(line, "FOREIGN KEY (`"); idx != -1 {
		start := idx + 14
		if end := strings.Index(line[start:], "`)"); end != -1 {
			fk.Column = line[start : start+end]
		}
	}

	// Buscar REFERENCES table (column)
	if idx := strings.Index(line, "REFERENCES `"); idx != -1 {
		start := idx + 12
		if end := strings.Index(line[start:], "`"); end != -1 {
			fk.ReferencedTable = line[start : start+end]
			// Buscar la columna referenciada
			refColStart := strings.Index(line[start+end:], "(`")
			if refColStart != -1 {
				refColStart += start + end + 2
				if refColEnd := strings.Index(line[refColStart:], "`)"); refColEnd != -1 {
					fk.ReferencedColumn = line[refColStart : refColStart+refColEnd]
				}
			}
		}
	}

	// Buscar ON DELETE/UPDATE
	if strings.Contains(line, "ON DELETE CASCADE") {
		fk.OnDelete = "CASCADE"
	} else if strings.Contains(line, "ON DELETE SET NULL") {
		fk.OnDelete = "SET NULL"
	}

	if strings.Contains(line, "ON UPDATE CASCADE") {
		fk.OnUpdate = "CASCADE"
	} else if strings.Contains(line, "ON UPDATE SET NULL") {
		fk.OnUpdate = "SET NULL"
	}

	// Verificar que tenemos información mínima
	if fk.Name != "" && fk.Column != "" && fk.ReferencedTable != "" && fk.ReferencedColumn != "" {
		return fk
	}

	return nil
}

func orderTablesByDependencies(db *sql.DB, tables []string) ([]string, error) {
	// Obtener todas las claves foráneas
	dependencies := make(map[string][]string)

	for _, table := range tables {
		fks, err := getTableForeignKeys(db, table)
		if err != nil {
			return nil, fmt.Errorf("error getting foreign keys for table %s: %v", table, err)
		}

		var deps []string
		for _, fk := range fks {
			if fk.ReferencedTable != table { // Evitar auto-referencias
				deps = append(deps, fk.ReferencedTable)
			}
		}
		dependencies[table] = deps
	}

	// Ordenamiento topológico simple
	var ordered []string
	visited := make(map[string]bool)

	var visit func(string)
	visit = func(table string) {
		if visited[table] {
			return
		}
		visited[table] = true

		for _, dep := range dependencies[table] {
			if contains(tables, dep) {
				visit(dep)
			}
		}
		ordered = append(ordered, table)
	}

	for _, table := range tables {
		visit(table)
	}

	return ordered, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateMigrationFile(tableInfo *TableInfo, timestamp time.Time) error {
	// Crear el directorio si no existe
	migrationDir := "database/migrations"
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return err
	}

	// Generar nombre del archivo
	timestampStr := timestamp.Format("2006_01_02_150405")
	fileName := fmt.Sprintf("%s_create_%s_table.go", timestampStr, tableInfo.Name)
	filePath := filepath.Join(migrationDir, fileName)

	// Verificar si el archivo ya existe
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Migration file already exists: %s\n", fileName)
		return nil
	}

	// Generar contenido del archivo
	content, err := generateMigrationContent(tableInfo, timestampStr)
	if err != nil {
		return err
	}

	// Escribir archivo
	return os.WriteFile(filePath, []byte(content), 0644)
}

func generateMigrationContent(tableInfo *TableInfo, timestamp string) (string, error) {
	tmpl := `package migrations

import (
	"database/sql"
	"web_utilidades/app/core/database"
)

type Create{{.StructName}}Table struct {
	database.BaseMigration
}

func NewCreate{{.StructName}}Table() *Create{{.StructName}}Table {
	return &Create{{.StructName}}Table{
		BaseMigration: database.BaseMigration{
			Name:      "create_{{.TableName}}_table",
			Timestamp: "{{.Timestamp}}",
		},
	}
}

func (m *Create{{.StructName}}Table) Up(db *sql.DB) error {
	query := ` + "`" + `
		CREATE TABLE {{.TableName}} ({{range $i, $col := .Columns}}{{if $i}},{{end}}
			{{$col.Name}} {{$col.SQLType}}{{$col.Constraints}}{{end}}{{if .Indexes}},{{range $i, $idx := .Indexes}}{{if $i}},{{end}}
			{{$idx.Definition}}{{end}}{{end}}{{if .ForeignKeys}},{{range $i, $fk := .ForeignKeys}}{{if $i}},{{end}}
			{{$fk.Definition}}{{end}}{{end}}
		)
	` + "`" + `
	_, err := db.Exec(query)
	return err
}

func (m *Create{{.StructName}}Table) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS {{.TableName}}")
	return err
}
`

	data := struct {
		StructName string
		TableName  string
		Timestamp  string
		Columns    []struct {
			Name        string
			SQLType     string
			Constraints string
		}
		Indexes []struct {
			Definition string
		}
		ForeignKeys []struct {
			Definition string
		}
	}{
		StructName: toPascalCase(tableInfo.Name),
		TableName:  tableInfo.Name,
		Timestamp:  timestamp,
	}

	// Procesar columnas
	for _, col := range tableInfo.Columns {
		sqlCol := struct {
			Name        string
			SQLType     string
			Constraints string
		}{
			Name:        col.Name,
			SQLType:     convertMySQLType(col.Type),
			Constraints: buildColumnConstraints(col),
		}
		data.Columns = append(data.Columns, sqlCol)
	}

	// Procesar índices
	for _, idx := range tableInfo.Indexes {
		indexDef := struct {
			Definition string
		}{
			Definition: buildIndexDefinition(idx),
		}
		data.Indexes = append(data.Indexes, indexDef)
	}

	// Procesar claves foráneas
	for _, fk := range tableInfo.ForeignKeys {
		fkDef := struct {
			Definition string
		}{
			Definition: buildForeignKeyDefinition(fk),
		}
		data.ForeignKeys = append(data.ForeignKeys, fkDef)
	}

	t, err := template.New("migration").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func convertMySQLType(mysqlType string) string {
	// Mantener los tipos de MySQL tal como están
	return strings.ToUpper(mysqlType)
}

func buildColumnConstraints(col ColumnInfo) string {
	var constraints []string

	// Primary key
	if col.Key == "PRI" {
		constraints = append(constraints, "PRIMARY KEY")
	}

	// Auto increment
	if col.Extra == "auto_increment" {
		constraints = append(constraints, "AUTO_INCREMENT")
	}

	// Unique
	if col.Key == "UNI" {
		constraints = append(constraints, "UNIQUE")
	}

	// Not null
	if col.Null == "NO" && col.Key != "PRI" {
		constraints = append(constraints, "NOT NULL")
	}

	// Default value
	if col.Default.Valid {
		if col.Default.String == "CURRENT_TIMESTAMP" {
			constraints = append(constraints, "DEFAULT CURRENT_TIMESTAMP")
		} else {
			constraints = append(constraints, fmt.Sprintf("DEFAULT '%s'", col.Default.String))
		}
	}

	if len(constraints) > 0 {
		return " " + strings.Join(constraints, " ")
	}
	return ""
}

func buildIndexDefinition(idx IndexInfo) string {
	columns := strings.Join(idx.Columns, ", ")
	if idx.Unique {
		return fmt.Sprintf("UNIQUE INDEX %s (%s)", idx.Name, columns)
	}
	return fmt.Sprintf("INDEX %s (%s)", idx.Name, columns)
}

func buildForeignKeyDefinition(fk ForeignKeyInfo) string {
	constraint := fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)",
		fk.Name, fk.Column, fk.ReferencedTable, fk.ReferencedColumn)

	if fk.OnDelete != "RESTRICT" {
		constraint += fmt.Sprintf(" ON DELETE %s", fk.OnDelete)
	}

	if fk.OnUpdate != "RESTRICT" {
		constraint += fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate)
	}

	return constraint
}
