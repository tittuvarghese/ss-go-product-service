package service

import (
	"fmt"
	"github.com/tittuvarghese/ss-go-product-service/core/database"
	"github.com/tittuvarghese/ss-go-product-service/models"
)

func CreateProduct(product models.Product, storage *database.RelationalDatabase) error {
	err := storage.Instance.Insert(&product)
	if err != nil {
		return err
	}
	return nil
}

func GetProduct(productId string, storage *database.RelationalDatabase) ([]models.Product, error) {
	var product []models.Product
	condition := map[string]interface{}{"id": productId}

	// Query the database with the given condition
	res, err := storage.Instance.QueryByCondition(&product, condition)
	if err != nil {
		return []models.Product{}, err
	}

	// Check if the result contains any products
	if len(res) <= 0 {
		return []models.Product{}, fmt.Errorf("product not found")
	}

	foundProduct, ok := res[0].(*[]models.Product) // Type assertion to pointer of models.Product
	if !ok {
		return []models.Product{}, fmt.Errorf("type assertion failed")
	}

	// Dereference the pointer to return the product
	return *foundProduct, nil
}

func GetProducts(storage *database.RelationalDatabase) (*[]models.Product, error) {
	var products []models.Product
	condition := map[string]interface{}{}

	// Pass a slice of Product to QueryByCondition
	res, err := storage.Instance.QueryByCondition(&products, condition)

	if err != nil {
		return nil, err
	}

	// Check if the result contains any product
	if len(res) == 0 {
		return nil, fmt.Errorf("products not found")
	}

	result, _ := res[0].(*[]models.Product)
	if len(*result) == 0 {
		return nil, fmt.Errorf("products not found")
	}

	return result, nil
}

func UpdateProduct(product models.Product, storage *database.RelationalDatabase) error {
	err := storage.Instance.Update(&product)
	if err != nil {
		return err
	}
	return nil
}
