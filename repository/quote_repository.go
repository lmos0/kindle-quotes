package repository

import (
	"database/sql"
	"quote-api/database"
	"quote-api/models"
	"time"
)

type QuoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{db: database.DB}
}

func (r *QuoteRepository) FindAll(limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT 
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) FindByID(id int64) (*models.Quote, error) {
	query := `
        SELECT 
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        WHERE q.id = ?
    `

	var quote models.Quote
	var book models.Book

	err := r.db.QueryRow(query, id).Scan(
		&quote.ID, &quote.BookID, &quote.Text, &quote.CreatedAt, &quote.UpdatedAt,
		&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
		&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	quote.Book = &book
	return &quote, nil
}

func (r *QuoteRepository) FindByBookID(bookID int64, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT 
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        WHERE q.book_id = ?
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, bookID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) FindByAuthorID(authorID int64, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT DISTINCT
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        INNER JOIN book_author ba ON b.id = ba.book_id
        WHERE ba.author_id = ?
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) FindByCategoryID(categoryID int, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT DISTINCT
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        INNER JOIN book_category bc ON b.id = bc.book_id
        WHERE bc.category_id = ?
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) FindByFilters(bookID *int64, authorID *int64, categoryID *int, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT DISTINCT
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
    `

	var conditions []string
	var args []interface{}

	if authorID != nil {
		query += " INNER JOIN book_author ba ON b.id = ba.book_id"
		conditions = append(conditions, "ba.author_id = ?")
		args = append(args, *authorID)
	}

	if categoryID != nil {
		query += " INNER JOIN book_category bc ON b.id = bc.book_id"
		conditions = append(conditions, "bc.category_id = ?")
		args = append(args, *categoryID)
	}

	if bookID != nil {
		conditions = append(conditions, "q.book_id = ?")
		args = append(args, *bookID)
	}

	if len(conditions) > 0 {
		query += " WHERE "
		for i, cond := range conditions {
			if i > 0 {
				query += " AND "
			}
			query += cond
		}
	}

	query += " ORDER BY q.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) FindRandom() (*models.Quote, error) {
	query := `
        SELECT 
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        ORDER BY RANDOM()
        LIMIT 1
    `

	var quote models.Quote
	var book models.Book

	err := r.db.QueryRow(query).Scan(
		&quote.ID, &quote.BookID, &quote.Text, &quote.CreatedAt, &quote.UpdatedAt,
		&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
		&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	quote.Book = &book
	return &quote, nil
}

func (r *QuoteRepository) Create(quote models.Quote) (*models.Quote, error) {
	query := `
        INSERT INTO quote (book_id, text, created_at, updated_at)
        VALUES (?, ?, ?, ?)
    `

	now := time.Now()
	result, err := r.db.Exec(query, quote.BookID, quote.Text, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *QuoteRepository) Update(id int64, quote models.Quote) (*models.Quote, error) {
	query := `
        UPDATE quote
        SET text = ?, updated_at = ?
        WHERE id = ?
    `

	now := time.Now()
	result, err := r.db.Exec(query, quote.Text, now, id)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindByID(id)
}

// Delete remove uma citação
func (r *QuoteRepository) Delete(id int64) error {
	query := "DELETE FROM quote WHERE id = ?"

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

func (r *QuoteRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM quote").Scan(&count)
	return count, err
}

func (r *QuoteRepository) CountByBookID(bookID int64) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM quote WHERE book_id = ?", bookID).Scan(&count)
	return count, err
}

func (r *QuoteRepository) CountByAuthorID(authorID int64) (int, error) {
	query := `
        SELECT COUNT(DISTINCT q.id)
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        INNER JOIN book_author ba ON b.id = ba.book_id
        WHERE ba.author_id = ?
    `

	var count int
	err := r.db.QueryRow(query, authorID).Scan(&count)
	return count, err
}

func (r *QuoteRepository) CountByCategoryID(categoryID int) (int, error) {
	query := `
        SELECT COUNT(DISTINCT q.id)
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        INNER JOIN book_category bc ON b.id = bc.book_id
        WHERE bc.category_id = ?
    `

	var count int
	err := r.db.QueryRow(query, categoryID).Scan(&count)
	return count, err
}

func (r *QuoteRepository) Search(searchTerm string, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT 
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        WHERE q.text LIKE ?
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.Query(query, "%"+searchTerm+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) SearchInBookAndAuthor(searchTerm string, limit, offset int) ([]models.Quote, error) {
	query := `
        SELECT DISTINCT
            q.id, q.book_id, q.text, q.created_at, q.updated_at,
            b.id, b.title, b.isbn, b.published_year, b.publisher, b.pages, b.created_at, b.updated_at
        FROM quote q
        INNER JOIN book b ON q.book_id = b.id
        LEFT JOIN book_author ba ON b.id = ba.book_id
        LEFT JOIN author a ON ba.author_id = a.id
        WHERE q.text LIKE ? OR b.title LIKE ? OR a.name LIKE ?
        ORDER BY q.created_at DESC
        LIMIT ? OFFSET ?
    `

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanQuotes(rows)
}

func (r *QuoteRepository) Exists(id int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM quote WHERE id = ?)"
	err := r.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}

func (r *QuoteRepository) BookExists(bookID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM book WHERE id = ?)"
	err := r.db.QueryRow(query, bookID).Scan(&exists)
	return exists, err
}

func (r *QuoteRepository) scanQuotes(rows *sql.Rows) ([]models.Quote, error) {
	var quotes []models.Quote

	for rows.Next() {
		var quote models.Quote
		var book models.Book

		err := rows.Scan(
			&quote.ID, &quote.BookID, &quote.Text, &quote.CreatedAt, &quote.UpdatedAt,
			&book.ID, &book.Title, &book.ISBN, &book.PublishedYear,
			&book.Publisher, &book.Pages, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		quote.Book = &book
		quotes = append(quotes, quote)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quotes, nil
}
