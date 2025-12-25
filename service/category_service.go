package service

import (
	"errors"
	"quote-api/models"
	"quote-api/repository"
	"strings"
)

type CategoryService struct {
	repo     repository.CategoryRepository
	bookRepo repository.BookRepository
}

func NewCategoryService(repo repository.CategoryRepository, bookRepo repository.BookRepository) *CategoryService {
	return &CategoryService{repo: repo, bookRepo: bookRepo}
}

func (s *CategoryService) getAll(limit, offset int) ([]models.Category, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	categories, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (s *CategoryService) getById(id int) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}

	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("Category not found")
	}

	return category, nil
}

func (s *CategoryService) Create(category models.Category) (*models.Category, error) {

	if err := s.validateCategory(category); err != nil {
		return nil, err
	}

	category.Name = strings.TrimSpace(category.Name)

	existing, err := s.repo.FindByID(category.ID)

	if err == nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("Category already exists")
	}

	created, err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *CategoryService) Update(id int, category models.Category) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}

	if err := s.validateCategory(category); err != nil {
		return nil, err
	}

	category.Name = strings.TrimSpace(category.Name)

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("Category not found")
	}

	duplicate, err := s.repo.FindByName(category.Name)
	if err != nil {
		return nil, err
	}

	if duplicate != nil && duplicate.ID != id {
		return nil, errors.New("Category already exists")
	}

	updated, err := s.repo.Update(id, category)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *CategoryService) Delete(id int) error {
	if id <= 0 {
		return errors.New("Invalid ID")
	}

	category, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if category == nil {
		return errors.New("Category not found")
	}

	books, err := s.bookRepo.FindByCategoryID(id, 1, 0)
	if err != nil {
		return err
	}
	if len(books) > 0 {
		return errors.New("Cannot delete category because there are still books in the category")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) Search(searchTerm string, limit, offset int) ([]models.Category, int, error) {
	if strings.TrimSpace(searchTerm) == "" {
		return s.getAll(limit, offset)
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	categories, err := s.repo.Search(searchTerm, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (s *CategoryService) GetBooksByCategory(categoryID int, limit, offset int) ([]models.Book, error) {
	if categoryID <= 0 {
		return nil, errors.New("ID de categoria inválido")
	}

	category, err := s.repo.FindByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	books, err := s.bookRepo.FindByCategoryID(categoryID, limit, offset)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *CategoryService) validateCategory(category models.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("Category name is required")
	}

	if len(category.Name) < 2 {
		return errors.New("Category name is too short")
	}

	if len(category.Name) > 100 {
		return errors.New("Category name is too long. Limit is 100 characters")
	}
	return nil
}
