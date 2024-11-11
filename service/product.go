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
		Height:                foundProduct.Height,
		Weight:                foundProduct.Weight,
		ShippingBasePrice:     foundProduct.ShippingBasePrice,
		BaseDeliveryTimelines: foundProduct.BaseDeliveryTimelines,
		SellerId:              foundProduct.SellerId,
	}

	return result, nil

}

func GetProducts(storage *database.RelationalDatabase) ([]models.Product, error) {
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

	// Initialize an empty result slice with the correct length
	result := make([]models.Product, 0, len(res))

	foundProduct, _ := res[0].(*[]models.Product)
	if len(*foundProduct) == 0 {
		return nil, fmt.Errorf("products not found")
	}
	// Iterate over the results and perform type assertion
	for _, product := range *foundProduct {
		// Construct the product data
		productData := models.Product{
			ID:                    product.ID,
			Name:                  product.Name,
			Quantity:              product.Quantity,
			Type:                  product.Type,
			Category:              product.Category,
			ImageUrls:             product.ImageUrls,
			Price:                 product.Price,
			Height:                product.Height,
			Weight:                product.Weight,
			ShippingBasePrice:     product.ShippingBasePrice,
			BaseDeliveryTimelines: product.BaseDeliveryTimelines,
			SellerId:              product.SellerId,
		}

		// Append the product to the result slice
		result = append(result, productData)

	}

	return result, nil
}
