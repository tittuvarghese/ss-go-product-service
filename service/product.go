package service

import (
	"fmt"
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

func GetProduct(productId string, storage *database.RelationalDatabase) (models.Product, error) {
	var product models.Product
	condition := map[string]interface{}{"id": productId}

	// Pass a slice of User to QueryByCondition
	res, err := storage.Instance.QueryByCondition(&product, condition)
	if err != nil {
		return models.Product{}, err
	}

	// Check if the result contains any user
	if len(res) == 0 {
		return models.Product{}, fmt.Errorf("product not found")
	}

	// Cast the result to the correct type (since QueryByCondition returns []interface{})
	foundProduct, ok := res[0].(*models.Product)
	if !ok {
		return models.Product{}, fmt.Errorf("type assertion failed")
	}

	result := models.Product{
		ID:                    foundProduct.ID,
		Name:                  foundProduct.Name,
		Quantity:              foundProduct.Quantity,
		Type:                  foundProduct.Type,
		Category:              foundProduct.Category,
		ImageUrls:             foundProduct.ImageUrls,
		Price:                 foundProduct.Price,
		Height:                foundProduct.Weight,
		Weight:                foundProduct.Weight,
		ShippingBasePrice:     foundProduct.ShippingBasePrice,
		BaseDeliveryTimelines: foundProduct.BaseDeliveryTimelines,
		SellerId:              foundProduct.SellerId,
	}

	return result, nil

}
