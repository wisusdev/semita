package seeders

import (
	"database/sql"
	"errors"
	"log"
	"semita/app/core/database"
	"semita/config"
)

// CategoriesSeeder seeder para categorías
type CategoriesSeeder struct {
	database.BaseSeeder
}

// NewCategoriesSeeder crea una nueva instancia del seeder
func NewCategoriesSeeder() *CategoriesSeeder {
	return &CategoriesSeeder{
		BaseSeeder: database.BaseSeeder{
			DB:   config.DatabaseConnect(),
			Name: "categories_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (cs *CategoriesSeeder) GetName() string {
	return cs.Name
}

// GetDependencies retorna las dependencias del seeder
func (cs *CategoriesSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// Seed ejecuta el seeding de categorías
func (cs *CategoriesSeeder) Seed() error {
	log.Println("Seeding categories...")

	categories := []struct {
		Name        string
		Description string
		Slug        string
	}{
		{
			Name:        "Tecnología",
			Description: "Artículos sobre tecnología, programación y desarrollo",
			Slug:        "tecnologia",
		},
		{
			Name:        "Ciencia",
			Description: "Contenido científico y descubrimientos",
			Slug:        "ciencia",
		},
		{
			Name:        "Deportes",
			Description: "Noticias y análisis deportivos",
			Slug:        "deportes",
		},
		{
			Name:        "Cultura",
			Description: "Arte, música, literatura y entretenimiento",
			Slug:        "cultura",
		},
		{
			Name:        "Negocios",
			Description: "Economía, finanzas y mundo empresarial",
			Slug:        "negocios",
		},
		{
			Name:        "Salud",
			Description: "Bienestar, medicina y vida saludable",
			Slug:        "salud",
		},
		{
			Name:        "Educación",
			Description: "Recursos educativos y aprendizaje",
			Slug:        "educacion",
		},
		{
			Name:        "Viajes",
			Description: "Destinos, experiencias y guías de viaje",
			Slug:        "viajes",
		},
	}

	for _, category := range categories {
		// Verificar si la categoría ya existe
		var existingID int
		checkQuery := `SELECT id FROM categories WHERE slug = ?`
		err := cs.DB.QueryRow(checkQuery, category.Slug).Scan(&existingID)

		if errors.Is(sql.ErrNoRows, err) {
			// No existe, crear nueva categoría
			insertQuery := `
				INSERT INTO categories (name, description, slug, created_at, updated_at) 
				VALUES (?, ?, ?, NOW(), NOW())`

			result, err := cs.DB.Exec(insertQuery, category.Name, category.Description, category.Slug)
			if err != nil {
				log.Printf("Error creating category '%s': %v", category.Name, err)
				continue
			}

			id, _ := result.LastInsertId()
			log.Printf("Created category: %s (ID: %d)", category.Name, id)
		} else if err != nil {
			log.Printf("Error checking existing category '%s': %v", category.Name, err)
			continue
		} else {
			log.Printf("Category '%s' already exists, skipping...", category.Name)
		}
	}

	log.Println("Categories seeding completed successfully!")
	return nil
}

// Rollback revierte el seeding de categorías
func (cs *CategoriesSeeder) Rollback() error {
	log.Println("Rolling back categories...")

	categorySlugs := []string{
		"tecnologia", "ciencia", "deportes", "cultura",
		"negocios", "salud", "educacion", "viajes",
	}

	for _, slug := range categorySlugs {
		query := `DELETE FROM categories WHERE slug = ?`
		result, err := cs.DB.Exec(query, slug)
		if err != nil {
			log.Printf("Error deleting category with slug '%s': %v", slug, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("Deleted category with slug: %s", slug)
		}
	}

	log.Println("Categories rollback completed successfully!")
	return nil
}
