package entities

import (
	"time"
)

type BasketProduct struct {
	BasketProductID string `pg:",pk" json:"basketproduct_id"`
	ProductCount    int    `json:"product_count"`

	ProductID string   ` json:"product_id"`
	Product   *Product `pg:"rel:has-one"`

	BasketID string  ` json:"basket_id"`
	Basket   *Basket `pg:", rel:has-one"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// status eklenecek
	// cost
}
