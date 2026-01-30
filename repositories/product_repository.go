package repositories

import (
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	products []models.Product
}

func NewProductRepository(products []models.Product) *ProductRepository {
	return &ProductRepository{
		products,
	}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	return repo.products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	product.ID = len(repo.products) + 1
	repo.products = append(repo.products, *product)

	return nil
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	for _, product := range repo.products {
		if product.ID == id {
			return &product, nil
		}
	}

	return nil, errors.New("Product not found")
}

func (repo *ProductRepository) Update(product *models.Product) error {
	for i := range repo.products {
		if repo.products[i].ID == product.ID {
			repo.products[i] = *product

			return nil
		}
	}

	return errors.New("Product not found")
}

func (repo *ProductRepository) Delete(id int) error {
	for i, product := range repo.products {
		if product.ID == id {
			repo.products = append(repo.products[:i], repo.products[i+1:]...)

			return nil
		}
	}

	return errors.New("Product not found")
}
