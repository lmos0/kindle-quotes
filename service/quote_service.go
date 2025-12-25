package service

import (
	"errors"
	"quote-api/models"
	"quote-api/repository"
	"strings"
)

type QuoteService struct {
	repo     *repository.QuoteRepository
	bookRepo *repository.BookRepository
}

func NewQuoteService(repo *repository.QuoteRepository, bookRepo *repository.BookRepository) *QuoteService {
	return &QuoteService{
		repo:     repo,
		bookRepo: bookRepo,
	}
}

func (s *QuoteService) GetAll(limit, offset int) ([]models.Quote, int, error) {

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) GetByID(id int64) (*models.Quote, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	quote, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, errors.New("citação não encontrada")
	}

	return quote, nil
}

func (s *QuoteService) GetByBookID(bookID int64, limit, offset int) ([]models.Quote, int, error) {
	if bookID <= 0 {
		return nil, 0, errors.New("ID de livro inválido")
	}

	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return nil, 0, err
	}
	if book == nil {
		return nil, 0, errors.New("livro não encontrado")
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.FindByBookID(bookID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByBookID(bookID)
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) GetByAuthorID(authorID int64, limit, offset int) ([]models.Quote, int, error) {
	if authorID <= 0 {
		return nil, 0, errors.New("ID de autor inválido")
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.FindByAuthorID(authorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByAuthorID(authorID)
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) GetByCategoryID(categoryID int, limit, offset int) ([]models.Quote, int, error) {
	if categoryID <= 0 {
		return nil, 0, errors.New("ID de categoria inválido")
	}

	// Validação de paginação
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.FindByCategoryID(categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByCategoryID(categoryID)
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) GetByFilters(bookID *int64, authorID *int64, categoryID *int, limit, offset int) ([]models.Quote, int, error) {

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.FindByFilters(bookID, authorID, categoryID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) GetRandom() (*models.Quote, error) {
	quote, err := s.repo.FindRandom()
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, errors.New("nenhuma citação disponível")
	}

	return quote, nil
}

func (s *QuoteService) Create(quote models.Quote) (*models.Quote, error) {

	if err := s.validateQuote(quote); err != nil {
		return nil, err
	}

	quote.Text = strings.TrimSpace(quote.Text)

	exists, err := s.repo.BookExists(quote.BookID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("livro não encontrado")
	}

	created, err := s.repo.Create(quote)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *QuoteService) Update(id int64, quote models.Quote) (*models.Quote, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	if err := s.validateQuote(quote); err != nil {
		return nil, err
	}

	quote.Text = strings.TrimSpace(quote.Text)

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("citação não encontrada")
	}

	updated, err := s.repo.Update(id, quote)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *QuoteService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}

	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("citação não encontrada")
	}
	
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *QuoteService) Search(searchTerm string, limit, offset int) ([]models.Quote, int, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return s.GetAll(limit, offset)
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.Search(searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return quotes, total, nil
}

func (s *QuoteService) SearchInBookAndAuthor(searchTerm string, limit, offset int) ([]models.Quote, int, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return s.GetAll(limit, offset)
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	quotes, err := s.repo.SearchInBookAndAuthor(searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total := len(quotes)

	return quotes, total, nil
}

func (s *QuoteService) validateQuote(quote models.Quote) error {
	if quote.BookID <= 0 {
		return errors.New("ID de livro inválido")
	}

	if strings.TrimSpace(quote.Text) == "" {
		return errors.New("texto da citação é obrigatório")
	}

	if len(quote.Text) < 3 {
		return errors.New("texto da citação deve ter pelo menos 3 caracteres")
	}

	if len(quote.Text) > 5000 {
		return errors.New("texto da citação deve ter no máximo 5000 caracteres")
	}

	return nil
}
