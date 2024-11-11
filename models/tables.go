package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;" json:"product_id"`
	Name                  string    `gorm:"type:varchar(255);not null" json:"name"`
	Quantity              int32     `gorm:"not null" json:"quantity"`
	Type                  string    `gorm:"type:varchar(20);not null" json:"type"`
	Category              string    `gorm:"type:varchar(100);not null" json:"category"`
	ImageURLs             string    `gorm:"type:json" json:"image_urls"`
	Price                 float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Width                 float64   `gorm:"type:decimal(5,2)" json:"width"`
	Height                float64   `gorm:"type:decimal(5,2)" json:"height"`
	Weight                float64   `gorm:"type:decimal(5,2)" json:"weight"`
	ShippingBasePrice     float64   `gorm:"type:decimal(10,2);not null" json:"shipping_base_price"`
	BaseDeliveryTimelines int32     `gorm:"not null" json:"base_delivery_timelines"`
	SellerID              uuid.UUID `gorm:"type:uuid;not null" json:"seller_id"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	product.ID = uuid.New()
	return nil
}
