package seeders

import (
	"database/sql"
	"log"
	"web_utilidades/app/core/database"
	"web_utilidades/config"
)

// PostsSeeder seeder para posts de ejemplo
type PostsSeeder struct {
	database.BaseSeeder
}

// NewPostsSeeder crea una nueva instancia del seeder
func NewPostsSeeder() *PostsSeeder {
	return &PostsSeeder{
		BaseSeeder: database.BaseSeeder{
			DB:   config.DatabaseConnect(),
			Name: "posts_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (ps *PostsSeeder) GetName() string {
	return ps.Name
}

// GetDependencies retorna las dependencias del seeder
func (ps *PostsSeeder) GetDependencies() []string {
	return []string{"users_seeder", "categories_seeder"} // Depende de usuarios y categorías
}

// Seed ejecuta el seeding de posts
func (ps *PostsSeeder) Seed() error {
	log.Println("Seeding posts...")

	posts := []struct {
		Title    string
		Content  string
		Slug     string
		Category string
		Author   string
	}{
		{
			Title:    "Introducción a Go: El lenguaje de programación del futuro",
			Content:  "Go, también conocido como Golang, es un lenguaje de programación desarrollado por Google. Su simplicidad, eficiencia y excelente soporte para concurrencia lo han convertido en una opción popular para el desarrollo de aplicaciones modernas...",
			Slug:     "introduccion-go-lenguaje-programacion",
			Category: "tecnologia",
			Author:   "editor@example.com",
		},
		{
			Title:    "Los beneficios del ejercicio regular para la salud mental",
			Content:  "El ejercicio regular no solo mejora nuestra condición física, sino que también tiene un impacto profundo en nuestra salud mental. Estudios recientes han demostrado que la actividad física puede reducir significativamente los síntomas de depresión y ansiedad...",
			Slug:     "beneficios-ejercicio-salud-mental",
			Category: "salud",
			Author:   "ana.martinez@example.com",
		},
		{
			Title:    "Descubrimiento de agua en exoplanetas: Un paso hacia la vida extraterrestre",
			Content:  "Los científicos han confirmado la presencia de vapor de agua en la atmósfera de varios exoplanetas. Este descubrimiento representa un hito importante en la búsqueda de vida fuera de nuestro sistema solar...",
			Slug:     "descubrimiento-agua-exoplanetas",
			Category: "ciencia",
			Author:   "editor@example.com",
		},
		{
			Title:    "El impacto de la inteligencia artificial en los negocios modernos",
			Content:  "La inteligencia artificial está transformando la manera en que las empresas operan y compiten en el mercado global. Desde la automatización de procesos hasta la personalización de experiencias del cliente...",
			Slug:     "impacto-ia-negocios-modernos",
			Category: "negocios",
			Author:   "admin@example.com",
		},
		{
			Title:    "Guía completa para viajar a Japón: Cultura, tradiciones y lugares imperdibles",
			Content:  "Japón es un destino fascinante que combina tradición milenaria con innovación tecnológica. Esta guía te ayudará a planificar tu viaje perfecto, desde los templos de Kioto hasta los rascacielos de Tokio...",
			Slug:     "guia-completa-viajar-japon",
			Category: "viajes",
			Author:   "user@example.com",
		},
		{
			Title:    "El renacimiento del arte digital: NFTs y nuevas formas de expresión",
			Content:  "Los tokens no fungibles (NFTs) han revolucionado el mundo del arte digital, creando nuevas oportunidades para artistas y coleccionistas. Exploramos cómo esta tecnología está cambiando el panorama cultural...",
			Slug:     "renacimiento-arte-digital-nfts",
			Category: "cultura",
			Author:   "ana.martinez@example.com",
		},
		{
			Title:    "Metodologías de aprendizaje efectivas en la era digital",
			Content:  "La educación ha evolucionado significativamente con la integración de tecnologías digitales. Exploramos las metodologías más efectivas para el aprendizaje en línea y híbrido...",
			Slug:     "metodologias-aprendizaje-era-digital",
			Category: "educacion",
			Author:   "editor@example.com",
		},
		{
			Title:    "Análisis del Mundial de Fútbol: Tendencias y estadísticas",
			Content:  "Un análisis profundo de las tendencias tácticas y estadísticas más relevantes del último Mundial de Fútbol. Desde la evolución del juego hasta el rendimiento de las selecciones participantes...",
			Slug:     "analisis-mundial-futbol-tendencias",
			Category: "deportes",
			Author:   "moderator@example.com",
		},
	}

	for _, post := range posts {
		// Verificar si el post ya existe
		var existingID int
		checkQuery := `SELECT id FROM posts WHERE slug = ?`
		err := ps.DB.QueryRow(checkQuery, post.Slug).Scan(&existingID)

		if err == sql.ErrNoRows {
			// Obtener el ID del autor
			authorID, err := ps.getUserIDByEmail(post.Author)
			if err != nil {
				log.Printf("Error getting author ID for '%s': %v", post.Author, err)
				continue
			}

			// No existe, crear nuevo post
			insertQuery := `
				INSERT INTO posts (title, content, slug, user_id, created_at, updated_at) 
				VALUES (?, ?, ?, ?, NOW(), NOW())`

			result, err := ps.DB.Exec(insertQuery, post.Title, post.Content, post.Slug, authorID)
			if err != nil {
				log.Printf("Error creating post '%s': %v", post.Title, err)
				continue
			}

			postID, _ := result.LastInsertId()
			log.Printf("Created post: %s (ID: %d)", post.Title, postID)
		} else if err != nil {
			log.Printf("Error checking existing post '%s': %v", post.Title, err)
			continue
		} else {
			log.Printf("Post '%s' already exists, skipping...", post.Title)
		}
	}

	log.Println("Posts seeding completed successfully!")
	return nil
}

// getUserIDByEmail obtiene el ID de un usuario por su email
func (ps *PostsSeeder) getUserIDByEmail(email string) (int, error) {
	var userID int
	query := `SELECT id FROM users WHERE email = ?`
	err := ps.DB.QueryRow(query, email).Scan(&userID)
	return userID, err
}

// Rollback revierte el seeding de posts
func (ps *PostsSeeder) Rollback() error {
	log.Println("Rolling back posts...")

	postSlugs := []string{
		"introduccion-go-lenguaje-programacion",
		"beneficios-ejercicio-salud-mental",
		"descubrimiento-agua-exoplanetas",
		"impacto-ia-negocios-modernos",
		"guia-completa-viajar-japon",
		"renacimiento-arte-digital-nfts",
		"metodologias-aprendizaje-era-digital",
		"analisis-mundial-futbol-tendencias",
	}

	for _, slug := range postSlugs {
		query := `DELETE FROM posts WHERE slug = ?`
		result, err := ps.DB.Exec(query, slug)
		if err != nil {
			log.Printf("Error deleting post with slug '%s': %v", slug, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("Deleted post with slug: %s", slug)
		}
	}

	log.Println("Posts rollback completed successfully!")
	return nil
}
