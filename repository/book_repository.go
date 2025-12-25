package repository

import (
	"database/sql"
	"quote-api/database"
	"quote-api/models"
	"time"
)

type BookRepository struct {
	db           *sql.DB
	authorRepo   *AuthorRepository
	categoryRepo *CategoryRepository
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		db:           database.DB,
		authorRepo:   NewAuthorRepository(),
		categoryRepo: NewCategoryRepository(),
	}
}

// FindAll lista todos os livros com autores e categorias
func (r *BookRepository) FindAll(limit, offset int) ([]models.Book, error) {
	query := `
        SELECT id, title, isbn, published_year, publisher, pages, created_at, updated_at
        FROM book
        ORDER BY title ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
			&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Carrega autores e categorias
		book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
		book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

		books = append(books, book)
	}

	return books, nil
}

// FindByID busca livro por ID com autores e categorias
func (r *BookRepository) FindByID(id int64) (*models.Book, error) {
	query := `
        SELECT id, title, isbn, published_year, publisher, pages, created_at, updated_at
        FROM book
        WHERE id = ?
    `

	var book models.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
		&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Carrega autores e categorias
	book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
	book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

	return &book, nil
}

// FindByISBN busca livro por ISBN
func (r *BookRepository) FindByISBN(isbn string) (*models.Book, error) {
	query := `
        SELECT id, title, isbn, published_year, publisher, pages, created_at, updated_at
        FROM book
        WHERE isbn = ?
    `

	var book models.Book
	err := r.db.QueryRow(query, isbn).Scan(
		&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
		&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
	book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

	return &book, nil
}

// FindByAuthorID busca livros de um autor
func (r *BookRepository) FindByAuthorID(authorID int64, limit, offset int) ([]models.Book, error) {
	query := `
        SELECT b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM book b
        INNER JOIN book_author ba ON b.id = ba.book_id
        WHERE ba.author_id = ?
        ORDER BY b.title ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
			&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
		book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

		books = append(books, book)
	}

	return books, nil
}

// FindByCategoryID busca livros de uma categoria
func (r *BookRepository) FindByCategoryID(categoryID int, limit, offset int) ([]models.Book, error) {
	query := `
        SELECT b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM book b
        INNER JOIN book_category bc ON b.id = bc.book_id
        WHERE bc.category_id = ?
        ORDER BY b.title ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
			&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
		book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

		books = append(books, book)
	}

	return books, nil
}

// Create cria livro com autores e categorias (usa transação)
func (r *BookRepository) Create(book models.Book, authorIDs []int64, categoryIDs []int) (*models.Book, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Insere livro
	query := `
        INSERT INTO book (title, isbn, published_year, publisher, pages, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	now := time.Now()
	result, err := tx.Exec(query, book.Title, book.ISBN, book.PublishedYear,
		book.Publisher, book.Pages, now, now)
	if err != nil {
		return nil, err
	}

	bookID, _ := result.LastInsertId()

	// 2. Associa autores (com ordem)
	for i, authorID := range authorIDs {
		_, err = tx.Exec(
			"INSERT INTO book_author (book_id, author_id, order) VALUES (?, ?, ?)",
			bookID, authorID, i+1,
		)
		if err != nil {
			return nil, err
		}
	}

	// 3. Associa categorias
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(
			"INSERT INTO book_category (book_id, category_id) VALUES (?, ?)",
			bookID, categoryID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.FindByID(bookID)
}

// Update atualiza livro (sem alterar autores/categorias)
func (r *BookRepository) Update(id int64, book models.Book) (*models.Book, error) {
	query := `
        UPDATE book
        SET title = ?, isbn = ?, published_year = ?, publisher = ?, pages = ?, updated_at = ?
        WHERE id = ?
    `

	now := time.Now()
	result, err := r.db.Exec(query, book.Title, book.ISBN, book.PublishedYear,
		book.Publisher, book.Pages, now, id)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(id)
}

// UpdateAuthors atualiza autores de um livro (usa transação)
func (r *BookRepository) UpdateAuthors(bookID int64, authorIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Remove todos os autores atuais
	_, err = tx.Exec("DELETE FROM book_author WHERE book_id = ?", bookID)
	if err != nil {
		return err
	}

	// 2. Adiciona novos autores
	for i, authorID := range authorIDs {
		_, err = tx.Exec(
			"INSERT INTO book_author (book_id, author_id, order) VALUES (?, ?, ?)",
			bookID, authorID, i+1,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateCategories atualiza categorias de um livro (usa transação)
func (r *BookRepository) UpdateCategories(bookID int64, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Remove todas as categorias atuais
	_, err = tx.Exec("DELETE FROM book_category WHERE book_id = ?", bookID)
	if err != nil {
		return err
	}

	// 2. Adiciona novas categorias
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec(
			"INSERT INTO book_category (book_id, category_id) VALUES (?, ?)",
			bookID, categoryID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete remove livro (CASCADE remove book_author e book_category)
func (r *BookRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM book WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Count conta total de livros
func (r *BookRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM book").Scan(&count)
	return count, err
}

// Search busca livros por título
func (r *BookRepository) Search(searchTerm string, limit, offset int) ([]models.Book, error) {
	query := `
        SELECT id, title, isbn, published_year, publisher, pages, created_at, updated_at
        FROM book
        WHERE title LIKE ?
        ORDER BY title ASC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
			&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		book.Authors, _ = r.authorRepo.FindByBookID(book.ID)
		book.Categories, _ = r.categoryRepo.FindByBookID(book.ID)

		books = append(books, book)
	}

	return books, nil
}
