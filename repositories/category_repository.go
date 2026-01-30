package repositories

import (
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	categories []models.Category
}

func NewCategoryRepository(categories []models.Category) *CategoryRepository {
	return &CategoryRepository{
		categories,
	}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	return repo.categories, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	category.ID = len(repo.categories) + 1
	repo.categories = append(repo.categories, *category)

	return nil
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	for _, category := range repo.categories {
		if category.ID == id {
			return &category, nil
		}
	}

	return nil, errors.New("Category not found")
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	for i := range repo.categories {
		if repo.categories[i].ID == category.ID {
			repo.categories[i] = *category

			return nil
		}
	}

	return errors.New("Category not found")
}

func (repo *CategoryRepository) Delete(id int) error {
	for i, category := range repo.categories {
		if category.ID == id {
			repo.categories = append(repo.categories[:i], repo.categories[i+1:]...)

			return nil
		}
	}

	return errors.New("Category not found")
}
