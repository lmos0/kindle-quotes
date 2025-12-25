package repository

import (
	"database/sql"
	"quote-api/database"
	"quote-api/models"
	"time"
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository() *AuthorRepository {
	return &AuthorRepository{db: database.DB}
}

func (r *AuthorRepository) FindAll(limit, offset int) ([]models.Author, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM author
        ORDER BY name ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var author models.Author
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (r *AuthorRepository) FindByID(id int64) (*models.Author, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM author
        WHERE id = ?
    `

	var author models.Author
	err := r.db.QueryRow(query, id).Scan(
		&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (r *AuthorRepository) FindByName(name string) (*models.Author, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM author
        WHERE LOWER(name) = LOWER(?)
    `

	var author models.Author
	err := r.db.QueryRow(query, name).Scan(
		&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (r *AuthorRepository) FindByBookID(bookID int64) ([]models.Author, error) {
	query := `
        SELECT a.id, a.name, a.created_at, a.updated_at
        FROM author a
        INNER JOIN book_author ba ON a.id = ba.author_id
        WHERE ba.book_id = ?
        ORDER BY ba.order ASC
    `

	rows, err := r.db.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var author models.Author
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (r *AuthorRepository) Create(author models.Author) (*models.Author, error) {
	query := `
        INSERT INTO author (name, created_at, updated_at)
        VALUES (?, ?, ?)
    `

	now := time.Now()
	result, err := r.db.Exec(query, author.Name, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *AuthorRepository) Update(id int64, author models.Author) (*models.Author, error) {
	query := `
        UPDATE author
        SET name = ?, updated_at = ?
        WHERE id = ?
    `

	now := time.Now()
	result, err := r.db.Exec(query, author.Name, now, id)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(id)
}

func (r *AuthorRepository) Delete(id int64) error {
	query := "DELETE FROM author WHERE id = ?"

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

func (r *AuthorRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM author").Scan(&count)
	return count, err
}

func (r *AuthorRepository) Search(searchTerm string, limit, offset int) ([]models.Author, error) {
	query := `
        SELECT id, name, created_at, updated_at
        FROM author
        WHERE name LIKE ?
        ORDER BY name ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var author models.Author
		err := rows.Scan(&author.ID, &author.Name, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
