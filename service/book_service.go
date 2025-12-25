package service

import (
	"errors"
	"quote-api/models"
	"quote-api/repository"
	"strings"
)

type BookService struct {
	repo         *repository.BookRepository
	authorRepo   *repository.AuthorRepository
	categoryRepo *repository.CategoryRepository
}

func NewBookService(
	repo *repository.BookRepository,
	authorRepo *repository.AuthorRepository,
	categoryRepo *repository.CategoryRepository,
) *BookService {
	return &BookService{
		repo:         repo,
		authorRepo:   authorRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *BookService) GetAll(limit, offset int) ([]models.Book, int, error) {

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *BookService) GetByID(id int64) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("livro não encontrado")
	}

	return book, nil
}

func (s *BookService) GetByISBN(isbn string) (*models.Book, error) {
	isbn = strings.TrimSpace(isbn)
	if isbn == "" {
		return nil, errors.New("ISBN inválido")
	}

	book, err := s.repo.FindByISBN(isbn)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("livro não encontrado")
	}

	return book, nil
}

func (s *BookService) Create(book models.Book, authorIDs []int64, categoryIDs []int) (*models.Book, error) {

	if err := s.validateBook(book); err != nil {
		return nil, err
	}

	if len(authorIDs) == 0 {
		return nil, errors.New("livro deve ter pelo menos um autor")
	}

	book.Title = strings.TrimSpace(book.Title)
	if book.Publisher != nil {
		publisher := strings.TrimSpace(*book.Publisher)
		book.Publisher = &publisher
	}

	if book.ISBN != nil && *book.ISBN != "" {
		existing, err := s.repo.FindByISBN(*book.ISBN)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("já existe um livro com este ISBN")
		}
	}

	for _, authorID := range authorIDs {
		author, err := s.authorRepo.FindByID(authorID)
		if err != nil {
			return nil, err
		}
		if author == nil {
			return nil, errors.New("um ou mais autores não encontrados")
		}
	}

	for _, categoryID := range categoryIDs {
		category, err := s.categoryRepo.FindByID(categoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("uma ou mais categorias não encontradas")
		}
	}

	created, err := s.repo.Create(book, authorIDs, categoryIDs)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *BookService) Update(id int64, book models.Book) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	if err := s.validateBook(book); err != nil {
		return nil, err
	}

	book.Title = strings.TrimSpace(book.Title)
	if book.Publisher != nil {
		publisher := strings.TrimSpace(*book.Publisher)
		book.Publisher = &publisher
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("livro não encontrado")
	}

	if book.ISBN != nil && *book.ISBN != "" {
		duplicate, err := s.repo.FindByISBN(*book.ISBN)
		if err != nil {
			return nil, err
		}
		if duplicate != nil && duplicate.ID != id {
			return nil, errors.New("já existe outro livro com este ISBN")
		}
	}

	updated, err := s.repo.Update(id, book)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *BookService) UpdateAuthors(bookID int64, authorIDs []int64) error {
	if bookID <= 0 {
		return errors.New("ID de livro inválido")
	}

	if len(authorIDs) == 0 {
		return errors.New("livro deve ter pelo menos um autor")
	}

	// Verifica se o livro existe
	book, err := s.repo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("livro não encontrado")
	}

	for _, authorID := range authorIDs {
		author, err := s.authorRepo.FindByID(authorID)
		if err != nil {
			return err
		}
		if author == nil {
			return errors.New("um ou mais autores não encontrados")
		}
	}

	err = s.repo.UpdateAuthors(bookID, authorIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookService) UpdateCategories(bookID int64, categoryIDs []int) error {
	if bookID <= 0 {
		return errors.New("ID de livro inválido")
	}

	book, err := s.repo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("livro não encontrado")
	}

	for _, categoryID := range categoryIDs {
		category, err := s.categoryRepo.FindByID(categoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return errors.New("uma ou mais categorias não encontradas")
		}
	}

	err = s.repo.UpdateCategories(bookID, categoryIDs)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}

	book, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("livro não encontrado")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookService) Search(searchTerm string, limit, offset int) ([]models.Book, int, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return s.GetAll(limit, offset)
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.repo.Search(searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *BookService) GetByAuthor(authorID int64, limit, offset int) ([]models.Book, int, error) {
	if authorID <= 0 {
		return nil, 0, errors.New("ID de autor inválido")
	}

	author, err := s.authorRepo.FindByID(authorID)
	if err != nil {
		return nil, 0, err
	}
	if author == nil {
		return nil, 0, errors.New("autor não encontrado")
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.repo.FindByAuthorID(authorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total := len(books)

	return books, total, nil
}

func (s *BookService) GetByCategory(categoryID int, limit, offset int) ([]models.Book, int, error) {
	if categoryID <= 0 {
		return nil, 0, errors.New("ID de categoria inválido")
	}

	category, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, 0, err
	}
	if category == nil {
		return nil, 0, errors.New("categoria não encontrada")
	}
	
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.repo.FindByCategoryID(categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total := len(books)

	return books, total, nil
}

func (s *BookService) validateBook(book models.Book) error {
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("título do livro é obrigatório")
	}

	if len(book.Title) < 1 {
		return errors.New("título do livro deve ter pelo menos 1 caractere")
	}

	if len(book.Title) > 500 {
		return errors.New("título do livro deve ter no máximo 500 caracteres")
	}

	if book.PublishedYear < 0 || book.PublishedYear > 9999 {
		return errors.New("ano de publicação inválido")
	}

	if book.Pages < 0 {
		return errors.New("número de páginas não pode ser negativo")
	}

	if book.ISBN != nil && *book.ISBN != "" {
		isbn := strings.ReplaceAll(*book.ISBN, "-", "")
		isbn = strings.ReplaceAll(isbn, " ", "")

		if len(isbn) != 10 && len(isbn) != 13 {
			return errors.New("ISBN deve ter 10 ou 13 dígitos")
		}
	}

	return nil
}
