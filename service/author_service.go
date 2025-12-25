package service

import (
	"errors"
	"quote-api/models"
	"quote-api/repository"
	"strings"
)

type AuthorService struct {
	repo     *repository.AuthorRepository
	bookRepo *repository.BookRepository
}

func NewAuthorService(repo *repository.AuthorRepository, bookRepo *repository.BookRepository) *AuthorService {
	return &AuthorService{
		repo:     repo,
		bookRepo: bookRepo,
	}
}

func (s *AuthorService) GetAll(limit, offset int) ([]models.Author, int, error) {

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	authors, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

func (s *AuthorService) GetByID(id int64) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	author, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("autor não encontrado")
	}

	return author, nil
}

func (s *AuthorService) Create(author models.Author) (*models.Author, error) {

	if err := s.validateAuthor(author); err != nil {
		return nil, err
	}

	author.Name = strings.TrimSpace(author.Name)

	existing, err := s.repo.FindByName(author.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("já existe um autor com este nome")
	}

	created, err := s.repo.Create(author)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *AuthorService) Update(id int64, author models.Author) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	if err := s.validateAuthor(author); err != nil {
		return nil, err
	}

	author.Name = strings.TrimSpace(author.Name)

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("autor não encontrado")
	}

	duplicate, err := s.repo.FindByName(author.Name)
	if err != nil {
		return nil, err
	}
	if duplicate != nil && duplicate.ID != id {
		return nil, errors.New("já existe outro autor com este nome")
	}

	updated, err := s.repo.Update(id, author)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *AuthorService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}

	author, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if author == nil {
		return errors.New("autor não encontrado")
	}

	books, err := s.bookRepo.FindByAuthorID(id, 1, 0)
	if err != nil {
		return err
	}
	if len(books) > 0 {
		return errors.New("não é possível deletar autor com livros associados")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthorService) Search(searchTerm string, limit, offset int) ([]models.Author, int, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return s.GetAll(limit, offset)
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	authors, err := s.repo.Search(searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Para simplificar, retorna o total geral (em produção, você pode querer contar apenas os resultados da busca)
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

func (s *AuthorService) GetBooksByAuthor(authorID int64, limit, offset int) ([]models.Book, error) {
	if authorID <= 0 {
		return nil, errors.New("ID de autor inválido")
	}

	author, err := s.repo.FindByID(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("autor não encontrado")
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.bookRepo.FindByAuthorID(authorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *AuthorService) validateAuthor(author models.Author) error {
	if strings.TrimSpace(author.Name) == "" {
		return errors.New("nome do autor é obrigatório")
	}

	if len(author.Name) < 2 {
		return errors.New("nome do autor deve ter pelo menos 2 caracteres")
	}

	if len(author.Name) > 200 {
		return errors.New("nome do autor deve ter no máximo 200 caracteres")
	}

	return nil
}
