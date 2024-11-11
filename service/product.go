package service

import (
	"github.com/tittuvarghese/product-service/core/database"
	"github.com/tittuvarghese/product-service/models"
)

func CreateProduct(product models.Product, storage *database.RelationalDatabase) error {
	err := storage.Instance.Insert(&product)
	if err != nil {
		return err
	}
	return nil
}
