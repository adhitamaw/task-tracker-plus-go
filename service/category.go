package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
)

type CategoryService interface {
	Store(category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	DeleteByName(name string) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
	GetListByUser(userID int) ([]model.Category, error)
}

type categoryService struct {
	categoryRepository repo.CategoryRepository
}

func NewCategoryService(categoryRepository repo.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository}
}

func (c *categoryService) Store(category *model.Category) error {
	err := c.categoryRepository.Store(category)
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryService) Update(id int, category model.Category) error {
	err := c.categoryRepository.Update(id, category)
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryService) Delete(id int) error {
	err := c.categoryRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryService) DeleteByName(name string) error {
	// First, get all categories to find the one with matching name
	categories, err := c.categoryRepository.GetList()
	if err != nil {
		return err
	}

	// Find category with matching name
	for _, category := range categories {
		if category.Name == name {
			return c.categoryRepository.Delete(category.ID)
		}
	}

	// If not found, that's ok - already "deleted"
	return nil
}

func (c *categoryService) GetByID(id int) (*model.Category, error) {
	category, err := c.categoryRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryService) GetList() ([]model.Category, error) {
	categories, err := c.categoryRepository.GetList()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryService) GetListByUser(userID int) ([]model.Category, error) {
	categories, err := c.categoryRepository.GetListByUser(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
