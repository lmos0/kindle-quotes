package repository

import (
	"database/sql"
	"quote-api/database"
	"quote-api/models"
	"time"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: database.DB}
}

// FindAll lista todas as categorias
func (r *CategoryRepository) FindAll(limit, offset int) ([]models.Category, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM category
        ORDER BY name ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindByID busca categoria por ID
func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM category
        WHERE id = ?
    `

	var category models.Category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByName busca categoria por nome
func (r *CategoryRepository) FindByName(name string) (*models.Category, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM category
        WHERE LOWER(name) = LOWER(?)
    `

	var category models.Category
	err := r.db.QueryRow(query, name).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// FindByBookID busca categorias de um livro específico
func (r *CategoryRepository) FindByBookID(bookID int64) ([]models.Category, error) {
	query := `
        SELECT c.id, c.name, c.created_at, c.updated_at
        FROM category c
        INNER JOIN book_category bc ON c.id = bc.category_id
        WHERE bc.book_id = ?
        ORDER BY c.name ASC
    `

	rows, err := r.db.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Create cria nova categoria
func (r *CategoryRepository) Create(category models.Category) (*models.Category, error) {
	query := `
        INSERT INTO category (name, created_at, updated_at)
        VALUES (?, ?, ?)
    `

	now := time.Now()
	result, err := r.db.Exec(query, category.Name, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(int(id))
}

// Update atualiza categoria
func (r *CategoryRepository) Update(id int, category models.Category) (*models.Category, error) {
	query := `
        UPDATE category
        SET name = ?, updated_at = ?
        WHERE id = ?
    `

	now := time.Now()
	result, err := r.db.Exec(query, category.Name, now, id)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(id)
}

// Delete remove categoria
func (r *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM category WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Count conta total de categorias
func (r *CategoryRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
	return count, err
}

// Search busca categorias por nome parcial
func (r *CategoryRepository) Search(searchTerm string, limit, offset int) ([]models.Category, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM category
        WHERE name LIKE ?
        ORDER BY name ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
